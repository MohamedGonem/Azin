#define AzAppVersion "0.2.2"

[Setup]
AppId={{CA1B358E-4F89-412E-B278-72C2F9B983BD}
AppName=Azin
AppVersion={#AzAppVersion}
AppPublisher=Azin Project
AppPublisherURL=https://github.com/azin-lang/Azin
AppSupportURL=https://github.com/azin-lang/Azin/issues
AppUpdatesURL=https://github.com/azin-lang/Azin/releases

DefaultDirName={autopf}\Azin
DefaultGroupName=Azin
DisableProgramGroupPage=yes

LicenseFile=..\..\LICENSE

VersionInfoVersion={#AzAppVersion}
VersionInfoCompany=Azin Project
VersionInfoDescription=Azin compiler installer
VersionInfoCopyright=Azin Project

OutputDir=..\installer
OutputBaseFilename=Azin-setup-{#AzAppVersion}

Compression=lzma2/ultra64
SolidCompression=yes

PrivilegesRequired=lowest
PrivilegesRequiredOverridesAllowed=dialog
ArchitecturesInstallIn64BitMode=x64compatible

ChangesEnvironment=yes
WizardStyle=modern
UsedUserAreasWarning=no

UninstallDisplayIcon={app}\azc.exe
CloseApplications=yes
RestartApplications=yes

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Tasks]
Name: "envpath"; Description: "Add Azin binary directory to PATH environment variable"; Flags: checkablealone

[Files]
Source: "..\..\build\azc.exe"; DestDir: "{app}"; DestName: "azc.exe"; Flags: ignoreversion

[Icons]
Name: "{group}\Azin"; Filename: "{app}\azc.exe"
Name: "{group}\Uninstall Azin"; Filename: "{uninstallexe}"

[Code]

const
  WM_SETTINGCHANGE = $001A;
  SMTO_ABORTIFHUNG = $0002;

function SendMessageTimeout(
  hWnd: LongInt;
  Msg: LongWord;
  wParam: LongInt;
  lParam: String;
  fuFlags: LongWord;
  uTimeout: LongWord;
  var lpdwResult: LongWord
): LongInt;
  external 'SendMessageTimeoutW@user32.dll stdcall';

function GetEnvRootKey: Integer;
begin
  if IsAdminInstallMode then
    Result := HKLM
  else
    Result := HKCU;
end;

function GetEnvSubKey: string;
begin
  if IsAdminInstallMode then
    Result := 'System\CurrentControlSet\Control\Session Manager\Environment'
  else
    Result := 'Environment';
end;

procedure RefreshEnvironment;
var
  MsgResult: LongWord;
begin
  SendMessageTimeout(
    $FFFF, { HWND_BROADCAST }
    WM_SETTINGCHANGE,
    0,
    'Environment',
    SMTO_ABORTIFHUNG,
    5000,
    MsgResult
  );
end;

function PathContains(Path, Dir: string): Boolean;
begin
  Result := Pos(';' + Lowercase(Dir) + ';', ';' + Lowercase(Path) + ';') > 0;
end;

procedure AddToPath(Dir: string);
var
  Path: string;
  RootKey: Integer;
  SubKey: string;
begin
  RootKey := GetEnvRootKey;
  SubKey := GetEnvSubKey;

  if not RegQueryStringValue(RootKey, SubKey, 'Path', Path) then
    Path := '';

  if not PathContains(Path, Dir) then
  begin
    if (Path <> '') and (Path[Length(Path)] <> ';') then
      Path := Path + ';';

    Path := Path + Dir;

    RegWriteExpandStringValue(RootKey, SubKey, 'Path', Path);
  end;
end;

procedure RemoveFromPath(const Dir: string);
var
  Path, NewPath, Entry: string;
  PosSep: Integer;
  RootKey: Integer;
  SubKey: string;
begin
  RootKey := GetEnvRootKey;
  SubKey := GetEnvSubKey;

  if not RegQueryStringValue(RootKey, SubKey, 'Path', Path) then
    Exit;

  NewPath := '';

  while Path <> '' do
  begin
    PosSep := Pos(';', Path);

    if PosSep = 0 then
    begin
      Entry := Path;
      Path := '';
    end
    else
    begin
      Entry := Copy(Path, 1, PosSep - 1);
      Delete(Path, 1, PosSep);
    end;

    if CompareText(Entry, Dir) <> 0 then
    begin
      if NewPath <> '' then
        NewPath := NewPath + ';';

      NewPath := NewPath + Entry;
    end;
  end;

  RegWriteExpandStringValue(RootKey, SubKey, 'Path', NewPath);
end;

procedure CurStepChanged(CurStep: TSetupStep);
begin
  if (CurStep = ssPostInstall) and WizardIsTaskSelected('envpath') then
  begin
    AddToPath(ExpandConstant('{app}'));
    RefreshEnvironment;
  end;
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
  if CurUninstallStep = usPostUninstall then
  begin
    RemoveFromPath(ExpandConstant('{app}'));
    RefreshEnvironment;
  end;
end;