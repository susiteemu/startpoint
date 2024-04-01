package profileui

type ProfileSelectedMsg struct {
	Profile Profile
}

type ProfileSelectCancelledMsg struct{}

type CreateProfileMsg struct{}

type DeleteProfileMsg struct {
	Profile Profile
}

type CopyProfileMsg struct {
	Profile Profile
}

type RenameProfileMsg struct {
	Profile Profile
}

type EditProfileMsg struct {
	Profile Profile
}
