// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs types_cgodefs.go

package win32

import "unsafe"

type (
	LONG_PTR  int32
	UINT_PTR  uint32
	ULONG_PTR uint32
	BYTE      uint8
	Int       int32
	LONG      int32
	DWORD     uint32
	UINT      uint32
	WORD      uint16
	USHORT    uint16
	LRESULT   int32
	TCHAR     uint16
	COLORREF  uint32

	HWND      unsafe.Pointer
	HDC       unsafe.Pointer
	HMENU     unsafe.Pointer
	HINSTANCE unsafe.Pointer
	HICON     unsafe.Pointer
	HCURSOR   unsafe.Pointer
	HBRUSH    unsafe.Pointer
	HRGN      unsafe.Pointer
	HMONITOR  unsafe.Pointer
	HHOOK     unsafe.Pointer
	HRAWINPUT unsafe.Pointer
)

const (
	ENUM_CURRENT_SETTINGS = 0xffffffff

	DISP_CHANGE_SUCCESSFUL  = 0x0
	DISP_CHANGE_BADDUALVIEW = -0x6
	DISP_CHANGE_BADFLAGS    = -0x4
	DISP_CHANGE_BADMODE     = -0x2
	DISP_CHANGE_BADPARAM    = -0x5
	DISP_CHANGE_FAILED      = -0x1
	DISP_CHANGE_NOTUPDATED  = -0x3
	DISP_CHANGE_RESTART     = 0x1

	CDS_TEST           = 0x2
	CDS_UPDATEREGISTRY = 0x1

	HORZSIZE      = 0x4
	VERTSIZE      = 0x6
	HORZRES       = 0x8
	VERTRES       = 0xa
	VREFRESH      = 0x74
	CM_GAMMA_RAMP = 0x2
)

type DISPLAY_DEVICE struct {
	Cb           uint32
	DeviceName   [32]uint16
	DeviceString [128]uint16
	StateFlags   uint32
	DeviceID     [128]uint16
	DeviceKey    [128]uint16
}
type DEVMODE struct {
	DmDeviceName       [32]uint16
	DmSpecVersion      uint16
	DmDriverVersion    uint16
	DmSize             uint16
	DmDriverExtra      uint16
	DmFields           uint32
	Anon0              [16]byte
	DmColor            int16
	DmDuplex           int16
	DmYResolution      int16
	DmTTOption         int16
	DmCollate          int16
	DmFormName         [32]uint16
	DmLogPixels        uint16
	DmBitsPerPel       uint32
	DmPelsWidth        uint32
	DmPelsHeight       uint32
	Anon1              [4]byte
	DmDisplayFrequency uint32
	DmICMMethod        uint32
	DmICMIntent        uint32
	DmMediaType        uint32
	DmDitherType       uint32
	DmReserved1        uint32
	DmReserved2        uint32
	DmPanningWidth     uint32
	DmPanningHeight    uint32
}
type POINTL struct {
	X int32
	Y int32
}
type GUID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]uint8
}
type POINT struct {
	X int32
	Y int32
}
type MINMAXINFO struct {
	PtReserved     POINT
	PtMaxSize      POINT
	PtMaxPosition  POINT
	PtMinTrackSize POINT
	PtMaxTrackSize POINT
}
type KBDLLHOOKSTRUCT struct {
	VkCode      uint32
	ScanCode    uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uint32
}
type PIXELFORMATDESCRIPTOR struct {
	NSize           uint16
	NVersion        uint16
	DwFlags         uint32
	IPixelType      uint8
	CColorBits      uint8
	CRedBits        uint8
	CRedShift       uint8
	CGreenBits      uint8
	CGreenShift     uint8
	CBlueBits       uint8
	CBlueShift      uint8
	CAlphaBits      uint8
	CAlphaShift     uint8
	CAccumBits      uint8
	CAccumRedBits   uint8
	CAccumGreenBits uint8
	CAccumBlueBits  uint8
	CAccumAlphaBits uint8
	CDepthBits      uint8
	CStencilBits    uint8
	CAuxBuffers     uint8
	ILayerType      uint8
	BReserved       uint8
	DwLayerMask     uint32
	DwVisibleMask   uint32
	DwDamageMask    uint32
}
type OSVERSIONINFOEX struct {
	DwOSVersionInfoSize uint32
	DwMajorVersion      uint32
	DwMinorVersion      uint32
	DwBuildNumber       uint32
	DwPlatformId        uint32
	SzCSDVersion        [128]uint16
	WServicePackMajor   uint16
	WServicePackMinor   uint16
	WSuiteMask          uint16
	WProductType        uint8
	WReserved           uint8
}
type RECT struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

const (
	RIDEV_INPUTSINK         = 0x00000100
	RID_INPUT               = 0x10000003
	RIM_TYPEMOUSE           = 0
	HID_USAGE_PAGE_GENERIC  = 0x01
	HID_USAGE_GENERIC_MOUSE = 0x02
	ICON_SMALL2             = 2
	WM_INPUT                = 0x00FF
	WM_XBUTTONDOWN          = 0x020B
	WM_XBUTTONUP            = 0x020C
	WM_XBUTTONDBLCLK        = 0x020D
	MK_CONTROL              = 0x0008
	MK_LBUTTON              = 0x0001
	MK_MBUTTON              = 0x0010
	MK_RBUTTON              = 0x0002
	MK_SHIFT                = 0x0004
	MK_XBUTTON1             = 0x0020
	MK_XBUTTON2             = 0x0040
	WM_MOUSEHWHEEL          = 0x020E

	DMDFO_DEFAULT = 0
	DMDFO_STRETCH = 1
	DMDFO_CENTER  = 2
)

type MONITORINFOEX struct {
	CbSize    uint32
	RcMonitor RECT
	RcWork    RECT
	DwFlags   uint32
	SzDevice  [32]uint16
}
type RAWINPUTHEADER struct {
	DwType  uint32
	DwSize  uint32
	HDevice *byte
	WParam  uint32
}
type RAWMOUSE struct {
	UsFlags            uint16
	Pad_cgo_0          [2]byte
	Anon0              [4]byte
	UlRawButtons       uint32
	LLastX             int32
	LLastY             int32
	UlExtraInformation uint32
}
type RAWINPUT struct {
	Header RAWINPUTHEADER
	Data   [24]byte
}
type VIDEOPARAMETERS struct {
	Guid                  GUID
	DwOffset              uint32
	DwCommand             uint32
	DwFlags               uint32
	DwMode                uint32
	DwTVStandard          uint32
	DwAvailableModes      uint32
	DwAvailableTVStandard uint32
	DwFlickerFilter       uint32
	DwOverScanX           uint32
	DwOverScanY           uint32
	DwMaxUnscaledX        uint32
	DwMaxUnscaledY        uint32
	DwPositionX           uint32
	DwPositionY           uint32
	DwBrightness          uint32
	DwContrast            uint32
	DwCPType              uint32
	DwCPCommand           uint32
	DwCPStandard          uint32
	DwCPKey               uint32
	APSTriggerBits        uint32
	BOEMCopyProtection    [256]uint8
}

type RAWINPUTDEVICE struct {
	UsUsagePage uint16
	UsUsage     uint16
	DwFlags     uint32
	HwndTarget  HWND
}
