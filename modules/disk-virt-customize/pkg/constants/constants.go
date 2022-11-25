package constants

// Exit codes
// reserve 0+ numbers for the exit code of the command
const (
	InvalidArguments              = -1 // same as go-arg invalid args exit
	PrepareGuestFSApplianceFailed = -2
	ExecuteFailed                 = -3
)

const (
	DiskImagePath                  = "/mnt/targetpvc/disk.img"
	GuestFSApplianceArchivePath    = "/usr/local/lib/guestfs/downloaded"
	VirtCustomizeCommandsFileName  = "virt_customize_commands"
	GuestFSApplianceArchivePathEnv = "LIBGUESTFS_APPLIANCE_ARCHIVE_PATH"
	TargetApplianceFolderEnv       = "LIBGUESTFS_APPLIANCE_TARGET_PATH"
	DiskVirtCustomizeHome          = "/home/disk-virt-customize"
)
