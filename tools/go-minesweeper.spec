%global debug_package %{nil}

Name:           go-minesweeper
Version:        0
Release:        %autorelease
Summary:        Minesweeper Game

%global package_id io.github.heathcliff26.%{name}

License:        Apache-2.0
URL:            https://github.com/heathcliff26/go-minesweeper
Source:         %{url}/archive/refs/tags/v%{version}.tar.gz

BuildRequires: golang >= 1.24
BuildRequires: gcc libXcursor-devel libXrandr-devel mesa-libGL-devel libXi-devel libXinerama-devel libXxf86vm-devel libxkbcommon-devel wayland-devel

%global _description %{expand:
This is an implementation of minesweeper in golang, made with the ui framework fyne.io.
I mainly created it because i was bored and wanted to create a gui app.}

%description %{_description}

%prep
%autosetup -n go-minesweeper-%{version} -p1

%build
make build

%install
install -D -m 755 bin/%{name} %{buildroot}/%{_bindir}/%{name}
install -D -m 644 packages/%{package_id}.desktop %{buildroot}/%{_datadir}/applications/%{package_id}.desktop
install -D -m 644 packages/%{package_id}.png %{buildroot}/%{_datadir}/icons/hicolor/512x512/apps/%{package_id}.png
install -D -m 644 packages/%{package_id}.svg %{buildroot}/%{_datadir}/icons/hicolor/scalable/apps/%{package_id}.svg
install -D -m 644 %{package_id}.metainfo.xml %{buildroot}/%{_datadir}/metainfo/%{package_id}.metainfo.xml

%files
%license LICENSE
%doc README.md
%{_bindir}/%{name}
%{_datadir}/applications/%{package_id}.desktop
%{_datadir}/icons/hicolor/512x512/apps/%{package_id}.png
%{_datadir}/icons/hicolor/scalable/apps/%{package_id}.svg
%{_datadir}/metainfo/%{package_id}.metainfo.xml

%changelog
%autochangelog
