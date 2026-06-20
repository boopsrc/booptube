; booptube Windows installer (Inno Setup)
; Build: iscc installer/windows/booptube.iss /DAppVersion=0.1.0

#ifndef AppVersion
#define AppVersion "dev"
#endif

#define AppName "booptube"
#define AppPublisher "booptube"
#define StagingDir "..\..\installer\staging"
#define OutputDir "..\..\.build"

[Setup]
AppId={{A1B2C3D4-E5F6-7890-ABCD-EF1234567890}
AppName={#AppName}
AppVersion={#AppVersion}
AppPublisher={#AppPublisher}
DefaultDirName={autopf}\{#AppName}
DefaultGroupName={#AppName}
OutputDir={#OutputDir}
OutputBaseFilename=booptube-{#AppVersion}-windows-amd64-setup
Compression=lzma2/ultra64
SolidCompression=yes
ArchitecturesInstallIn64BitMode=x64compatible
PrivilegesRequired=lowest
WizardStyle=modern
DisableProgramGroupPage=no

[Languages]
Name: "brazilianportuguese"; MessagesFile: "compiler:Languages\BrazilianPortuguese.isl"
Name: "english"; MessagesFile: "compiler:Default.isl"

[Tasks]
Name: "desktopicon"; Description: "Atalho na area de trabalho (GUI)"; GroupDescription: "Atalhos:"; Flags: unchecked

[Components]
Name: "gui"; Description: "booptube-gui (interface grafica)"; Types: full custom; Flags: fixed
Name: "cli"; Description: "booptube (linha de comando)"; Types: full custom

[Files]
Source: "{#StagingDir}\booptube-gui.exe"; DestDir: "{app}"; Components: gui; Flags: ignoreversion
Source: "{#StagingDir}\booptube.exe"; DestDir: "{app}"; Components: cli; Flags: ignoreversion
Source: "{#StagingDir}\tools\*"; DestDir: "{app}\tools"; Components: gui cli; Flags: ignoreversion recursesubdirs
Source: "{#StagingDir}\README.md"; DestDir: "{app}"; Flags: ignoreversion skipifsourcedoesntexist
Source: "{#StagingDir}\VERSION.txt"; DestDir: "{app}"; Flags: ignoreversion skipifsourcedoesntexist

[Icons]
Name: "{group}\booptube-gui"; Filename: "{app}\booptube-gui.exe"; Components: gui
Name: "{group}\booptube (CLI)"; Filename: "{app}\booptube.exe"; Components: cli
Name: "{autodesktop}\booptube-gui"; Filename: "{app}\booptube-gui.exe"; Tasks: desktopicon; Components: gui

[Run]
Filename: "{app}\booptube-gui.exe"; Description: "Abrir booptube-gui"; Flags: nowait postinstall skipifsilent; Components: gui

[UninstallDelete]
Type: filesandordirs; Name: "{app}\tools"
