package attacks

import "errors"

var (
	// ErrMinimumMaxTime is returned when the duration provided is less than the minimum time required.
	ErrMinimumMaxTime error = errors.New("duration must be greater than or equal to minimum time")

	// ErrMaxtimeGreaterThanSelf is returned when the duration provided is greater than the maximum allowed time.
	ErrMaxtimeGreaterThanSelf error = errors.New("duration presented is greater than your maxtime")

	// ErrMaxtimeGreaterThanOverride is returned when the duration provided is greater than the overridden maximum time.
	ErrMaxtimeGreaterThanOverride error = errors.New("duration presented is greater than the method override")

	// ErrMaxtimeGreaterThanEnhanced is returned when the duration provided is greater than the maximum time set for the method.
	ErrMaxtimeGreaterThanEnhanced error = errors.New("duration presented is greater than the method maxtime")

	// ErrConcurrentLimitReached is returned when the concurrent limit for a user has been reached.
	ErrConcurrentLimitReached error = errors.New("user concurrent limit reached")

	// ErrGroupConcurrentLimitReached is returned when the group limit for a method has been reached.
	ErrGroupConcurrentLimitReached error = errors.New("group concurrent limit reached")

	// ErrCooldownPeriodEngaged is whenever they're inside a cooldown period
	ErrCooldownPeriodEngaged error = errors.New("cooldown period engaged")

	// ErrInvalidTarget is whenever the target they've provided is invalid
	ErrInvalidTarget error = errors.New("invalid target provided")

	// ErrUserBlacklisted is when the target is blacklisted by a user, has an expiry.
	ErrUserBlacklisted error = errors.New("user blacklisted")

	// ErrSystemBlacklisted is when the target is blacklisted by the toml file
	ErrSystemBlacklisted error = errors.New("system blacklisted")

	// ErrUnauthorizedSender is triggered when the sender is not authorized
	ErrUnauthorizedSender error = errors.New("sender is not authorized")

	// ErrPowersaving is triggered when the target is already being attacked
	ErrPowersaving error = errors.New("powersaving enabled")

	// ErrOngoingMethodCap is when the method has enough ongoing attacks
	ErrOngoingMethodCap error = errors.New("method ongoing limit reached")

	// ErrGlobalConns is when there aren't any slots
	ErrGlobalConns error = errors.New("no open global slots")

	// ErrGlobalMaxtime is whenever the maxtime requested is above the global max
	ErrGlobalMaxtime error = errors.New("maxtime breaches the global limit")

	// ErrGlobalCooldown is shown when the CNC is in a global cooldown state
	ErrGlobalCooldown error = errors.New("global cooldown engaged")

	// ErrGroupCooldown is when the group is in cooldown
	ErrGroupCooldown error = errors.New("group cooldown")

	// ErrMethodDisabled is shown when the method is disabled
	ErrMethodDisabled error = errors.New("method disabled")

	// ErrMaxAttacksReached is shown when the daily attacks limit is reached
	ErrMaxAttacksToday error = errors.New("your daily max attacks has been reached")

	// ErrAttacksDisabled is when the entire CNCs attacks are disabled
	ErrAttacksDisabled error = errors.New("attacks are now disabled")
)