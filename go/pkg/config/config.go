package config

const (
	DataDirectory      = "MHWArchiveManager_Data"
	IndexFile          = "MHWArchiveManager_Index.txt"
	LaunchFile         = "MHWArchiveManager_Launcher.txt"
	ExclusionFile      = "MHWArchiveManager_Exclusions.txt"
	SavedIndexPathFile = "MHWArchiveManager_IndexSave.txt"
	ProfileListFile    = "MHWArchiveManager_ProfileList.txt"
	DefaultProfileName = "DefaultProfile"
	TempDirectory      = "Temp"
	OutputDirectory    = "nativePC"
)

var ValidFileTypes = []string{".zip", ".7z", ".rar"}
