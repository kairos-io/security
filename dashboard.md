# Kairos Security Dashboard

_Updated 2026-06-22._

> No security findings to triage this run.

## 🔥 Focus now

_Nothing flagged._

## 🌊 Waterfall fronts

_None._

## 📦 Per-repo findings

| Repo | Critical | High | Medium | Low | Total | Status |
|---|---|---|---|---|---|---|
| kairos-io/AuroraBoot | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/cluster-api-provider-kairos | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/entangle | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/entangle-proxy | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/go-nodepair | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/go-ukify | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/hadron | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/immucore | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/kairos | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/kairos-agent | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/kairos-init | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/kairos-installer | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/kairos-lab | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/kairos-must-burn | 0 | 0 | 0 | 0 | 0 | ⚠️ errors |
| kairos-io/kairos-operator | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/kairos-sdk | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/kcrypt | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/kcrypt-discovery-challenger | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/netboot | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/provider-kairos | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/provider-kubernetes | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/simple-mdns-server | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/tpm-helpers | 0 | 0 | 0 | 0 | 0 | clean |
| mauromorales/xpasswd | 0 | 0 | 0 | 0 | 0 | clean |
| mudler/edgevpn | 0 | 0 | 0 | 0 | 0 | clean |
| mudler/entities | 0 | 0 | 0 | 0 | 0 | clean |
| mudler/go-pluggable | 0 | 0 | 0 | 0 | 0 | clean |
| mudler/yip | 0 | 0 | 0 | 0 | 0 | clean |

## ⚠️ 1 collection errors

- `kairos-io/kairos-must-burn` / sourceCVE: govulncheck: exit status 1: govulncheck: loading packages: 
There are errors with the provided package patterns:

-: # github.com/diamondburned/gotk4/pkg/core/gbox
# [pkg-config --cflags  -- glib-2.0]
Package glib-2.0 was not found in the pkg-config search path.
Perhaps you should add the directory containing `glib-2.0.pc'
to the PKG_CONFIG_PATH environment variable
Package 'glib-2.0', required by 'virtual:world', not found
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/core/gbox/box.go:7:8: could not import C (no metadata for C)
-: # github.com/diamondburned/gotk4/pkg/core/intern
# [pkg-config --cflags  -- gobject-2.0 gobject-2.0]
Package gobject-2.0 was not found in the pkg-config search path.
Perhaps you should add the directory containing `gobject-2.0.pc'
to the PKG_CONFIG_PATH environment variable
Package 'gobject-2.0', required by 'virtual:world', not found
Package 'gobject-2.0', required by 'virtual:world', not found
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/core/intern/intern.go:6:8: could not import C (no metadata for C)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/core/glib/connect.go:8:8: could not import C (no metadata for C)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/core/gcancel/gcancel.go:7:8: could not import C (no metadata for C)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/core/gerror/gerror.go:9:8: could not import C (no metadata for C)
-: # github.com/diamondburned/gotk4/pkg/core/gextras
# [pkg-config --cflags  -- glib-2.0]
Package glib-2.0 was not found in the pkg-config search path.
Perhaps you should add the directory containing `glib-2.0.pc'
to the PKG_CONFIG_PATH environment variable
Package 'glib-2.0', required by 'virtual:world', not found
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/core/gextras/gextras.go:7:8: could not import C (no metadata for C)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:36:8: could not import C (no metadata for C)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:1208:23: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type NormalizeMode
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:2053:26: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type UnicodeBreakType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:2507:23: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type UnicodeScript
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:2990:21: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type UnicodeType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3388:10: invalid operation: f == 0 (mismatched types FileSetContentsFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3395:11: invalid operation: f != 0 (mismatched types FileSetContentsFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3396:16: invalid operation: f - 1 (mismatched types FileSetContentsFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3444:10: invalid operation: f == 0 (mismatched types FileTest and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3451:11: invalid operation: f != 0 (mismatched types FileTest and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3452:16: invalid operation: f - 1 (mismatched types FileTest and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3511:10: invalid operation: f == 0 (mismatched types FormatSizeFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3518:11: invalid operation: f != 0 (mismatched types FormatSizeFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3519:16: invalid operation: f - 1 (mismatched types FormatSizeFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3571:21: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type IOCondition
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3576:10: invalid operation: i == 0 (mismatched types IOCondition and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3583:11: invalid operation: i != 0 (mismatched types IOCondition and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3584:16: invalid operation: i - 1 (mismatched types IOCondition and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3655:10: invalid operation: i == 0 (mismatched types IOFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3662:11: invalid operation: i != 0 (mismatched types IOFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3663:16: invalid operation: i - 1 (mismatched types IOFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3716:10: invalid operation: k == 0 (mismatched types KeyFileFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3723:11: invalid operation: k != 0 (mismatched types KeyFileFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3724:16: invalid operation: k - 1 (mismatched types KeyFileFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3779:10: invalid operation: l == 0 (mismatched types LogLevelFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3786:11: invalid operation: l != 0 (mismatched types LogLevelFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3787:16: invalid operation: l - 1 (mismatched types LogLevelFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3838:10: invalid operation: m == 0 (mismatched types MainContextFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3845:11: invalid operation: m != 0 (mismatched types MainContextFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3846:16: invalid operation: m - 1 (mismatched types MainContextFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3907:10: invalid operation: m == 0 (mismatched types MarkupCollectType and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3914:11: invalid operation: m != 0 (mismatched types MarkupCollectType and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3915:16: invalid operation: m - 1 (mismatched types MarkupCollectType and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3972:10: invalid operation: m == 0 (mismatched types MarkupParseFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3979:11: invalid operation: m != 0 (mismatched types MarkupParseFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:3980:16: invalid operation: m - 1 (mismatched types MarkupParseFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:4047:10: invalid operation: o == 0 (mismatched types OptionFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:4054:11: invalid operation: o != 0 (mismatched types OptionFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:4055:16: invalid operation: o - 1 (mismatched types OptionFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:4187:10: invalid operation: r == 0 (mismatched types RegexCompileFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:4194:11: invalid operation: r != 0 (mismatched types RegexCompileFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:4195:16: invalid operation: r - 1 (mismatched types RegexCompileFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:4333:10: invalid operation: r == 0 (mismatched types RegexMatchFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:4340:11: invalid operation: r != 0 (mismatched types RegexMatchFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:4341:16: invalid operation: r - 1 (mismatched types RegexMatchFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:4438:10: invalid operation: s == 0 (mismatched types SpawnFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:4445:11: invalid operation: s != 0 (mismatched types SpawnFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:4446:16: invalid operation: s - 1 (mismatched types SpawnFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:4514:10: invalid operation: t == 0 (mismatched types TraverseFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:4521:11: invalid operation: t != 0 (mismatched types TraverseFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:4522:16: invalid operation: t - 1 (mismatched types TraverseFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:4597:10: invalid operation: u == 0 (mismatched types URIFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:4604:11: invalid operation: u != 0 (mismatched types URIFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:4605:16: invalid operation: u - 1 (mismatched types URIFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:4667:10: invalid operation: u == 0 (mismatched types URIHideFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:4674:11: invalid operation: u != 0 (mismatched types URIHideFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:4675:16: invalid operation: u - 1 (mismatched types URIHideFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:4724:10: invalid operation: u == 0 (mismatched types URIParamsFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:4731:11: invalid operation: u != 0 (mismatched types URIParamsFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/glib/v2/glib.go:4732:16: invalid operation: u - 1 (mismatched types URIParamsFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:1828:8: could not import C (no metadata for C)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:3269:17: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type BusType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:3303:25: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type ConverterResult
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:3352:25: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type CredentialsType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:3512:19: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type DBusError
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:3945:30: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type DBusMessageByteOrder
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:3990:32: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type DBusMessageHeaderField
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:4038:25: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type DBusMessageType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:4074:29: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type DataStreamByteOrder
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:4110:31: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type DataStreamNewlineType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:4152:28: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type DriveStartStopType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:4191:22: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type EmblemOrigin
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:4224:29: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type FileAttributeStatus
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:4268:27: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type FileAttributeType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:4335:26: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type FileMonitorEvent
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:4400:18: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type FileType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:4442:31: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type FilesystemPreviewType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:4597:21: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type IOErrorEnum
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:4719:28: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type IOModuleScopeFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:4760:35: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type MemoryMonitorWarningLevel
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:4792:30: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type MountOperationResult
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:4830:29: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type NetworkConnectivity
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:4872:30: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type NotificationPriority
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:4907:22: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type PasswordSave
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:4946:24: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type PollableReturn
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:4978:23: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type ResolverError
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5058:28: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type ResolverRecordType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5091:23: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type ResourceError
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5153:27: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type SocketClientEvent
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5198:22: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type SocketFamily
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5236:29: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type SocketListenerEvent
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5278:24: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type SocketProtocol
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5317:20: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type SocketType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5349:31: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type TLSAuthenticationMode
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5377:36: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type TLSCertificateRequestFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5417:32: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type TLSChannelBindingError
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5474:31: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type TLSChannelBindingType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5506:32: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type TLSDatabaseLookupFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5556:18: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type TLSError
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5619:30: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type TLSInteractionResult
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5674:28: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type TLSProtocolVersion
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5719:28: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type TLSRehandshakeMode
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5750:30: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type ZlibCompressorFormat
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5783:28: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type AppInfoCreateFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5788:10: invalid operation: a == 0 (mismatched types AppInfoCreateFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5795:11: invalid operation: a != 0 (mismatched types AppInfoCreateFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5796:16: invalid operation: a - 1 (mismatched types AppInfoCreateFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5875:26: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type ApplicationFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5880:10: invalid operation: a == 0 (mismatched types ApplicationFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5887:11: invalid operation: a != 0 (mismatched types ApplicationFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5888:16: invalid operation: a - 1 (mismatched types ApplicationFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5947:26: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type AskPasswordFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5952:10: invalid operation: a == 0 (mismatched types AskPasswordFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5959:11: invalid operation: a != 0 (mismatched types AskPasswordFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:5960:16: invalid operation: a - 1 (mismatched types AskPasswordFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6011:27: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type BusNameOwnerFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6016:10: invalid operation: b == 0 (mismatched types BusNameOwnerFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6023:11: invalid operation: b != 0 (mismatched types BusNameOwnerFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6024:16: invalid operation: b - 1 (mismatched types BusNameOwnerFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6063:29: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type BusNameWatcherFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6068:10: invalid operation: b == 0 (mismatched types BusNameWatcherFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6075:11: invalid operation: b != 0 (mismatched types BusNameWatcherFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6076:16: invalid operation: b - 1 (mismatched types BusNameWatcherFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6112:24: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type ConverterFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6117:10: invalid operation: c == 0 (mismatched types ConverterFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6124:11: invalid operation: c != 0 (mismatched types ConverterFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6125:16: invalid operation: c - 1 (mismatched types ConverterFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6165:23: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type DBusCallFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6170:10: invalid operation: d == 0 (mismatched types DBusCallFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6177:11: invalid operation: d != 0 (mismatched types DBusCallFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6178:16: invalid operation: d - 1 (mismatched types DBusCallFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6215:29: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type DBusCapabilityFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6220:10: invalid operation: d == 0 (mismatched types DBusCapabilityFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6227:11: invalid operation: d != 0 (mismatched types DBusCapabilityFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6228:16: invalid operation: d - 1 (mismatched types DBusCapabilityFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6289:29: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type DBusConnectionFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6294:10: invalid operation: d == 0 (mismatched types DBusConnectionFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6301:11: invalid operation: d != 0 (mismatched types DBusConnectionFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6302:16: invalid operation: d - 1 (mismatched types DBusConnectionFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6353:36: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type DBusInterfaceSkeletonFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6358:10: invalid operation: d == 0 (mismatched types DBusInterfaceSkeletonFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6365:11: invalid operation: d != 0 (mismatched types DBusInterfaceSkeletonFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6366:16: invalid operation: d - 1 (mismatched types DBusInterfaceSkeletonFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6407:26: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type DBusMessageFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6412:10: invalid operation: d == 0 (mismatched types DBusMessageFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6419:11: invalid operation: d != 0 (mismatched types DBusMessageFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6420:16: invalid operation: d - 1 (mismatched types DBusMessageFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6462:38: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type DBusObjectManagerClientFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6467:10: invalid operation: d == 0 (mismatched types DBusObjectManagerClientFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6474:11: invalid operation: d != 0 (mismatched types DBusObjectManagerClientFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6475:16: invalid operation: d - 1 (mismatched types DBusObjectManagerClientFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6512:31: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type DBusPropertyInfoFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6517:10: invalid operation: d == 0 (mismatched types DBusPropertyInfoFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6524:11: invalid operation: d != 0 (mismatched types DBusPropertyInfoFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6525:16: invalid operation: d - 1 (mismatched types DBusPropertyInfoFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6589:24: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type DBusProxyFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6594:10: invalid operation: d == 0 (mismatched types DBusProxyFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6601:11: invalid operation: d != 0 (mismatched types DBusProxyFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6602:16: invalid operation: d - 1 (mismatched types DBusProxyFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6647:30: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type DBusSendMessageFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6652:10: invalid operation: d == 0 (mismatched types DBusSendMessageFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6659:11: invalid operation: d != 0 (mismatched types DBusSendMessageFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6660:16: invalid operation: d - 1 (mismatched types DBusSendMessageFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6702:25: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type DBusServerFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6707:10: invalid operation: d == 0 (mismatched types DBusServerFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6714:11: invalid operation: d != 0 (mismatched types DBusServerFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6715:16: invalid operation: d - 1 (mismatched types DBusServerFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6763:25: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type DBusSignalFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6768:10: invalid operation: d == 0 (mismatched types DBusSignalFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6775:11: invalid operation: d != 0 (mismatched types DBusSignalFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6776:16: invalid operation: d - 1 (mismatched types DBusSignalFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6816:26: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type DBusSubtreeFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6821:10: invalid operation: d == 0 (mismatched types DBusSubtreeFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6828:11: invalid operation: d != 0 (mismatched types DBusSubtreeFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6829:16: invalid operation: d - 1 (mismatched types DBusSubtreeFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6861:25: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type DriveStartFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6866:10: invalid operation: d == 0 (mismatched types DriveStartFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6873:11: invalid operation: d != 0 (mismatched types DriveStartFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6874:16: invalid operation: d - 1 (mismatched types DriveStartFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6910:32: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type FileAttributeInfoFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6915:10: invalid operation: f == 0 (mismatched types FileAttributeInfoFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6922:11: invalid operation: f != 0 (mismatched types FileAttributeInfoFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6923:16: invalid operation: f - 1 (mismatched types FileAttributeInfoFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6975:23: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type FileCopyFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6980:10: invalid operation: f == 0 (mismatched types FileCopyFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6987:11: invalid operation: f != 0 (mismatched types FileCopyFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:6988:16: invalid operation: f - 1 (mismatched types FileCopyFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7044:25: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type FileCreateFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7049:10: invalid operation: f == 0 (mismatched types FileCreateFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7056:11: invalid operation: f != 0 (mismatched types FileCreateFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7057:16: invalid operation: f - 1 (mismatched types FileCreateFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7106:26: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type FileMeasureFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7111:10: invalid operation: f == 0 (mismatched types FileMeasureFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7118:11: invalid operation: f != 0 (mismatched types FileMeasureFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7119:16: invalid operation: f - 1 (mismatched types FileMeasureFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7171:26: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type FileMonitorFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7176:10: invalid operation: f == 0 (mismatched types FileMonitorFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7183:11: invalid operation: f != 0 (mismatched types FileMonitorFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7184:16: invalid operation: f - 1 (mismatched types FileMonitorFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7224:28: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type FileQueryInfoFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7229:10: invalid operation: f == 0 (mismatched types FileQueryInfoFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7236:11: invalid operation: f != 0 (mismatched types FileQueryInfoFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7237:16: invalid operation: f - 1 (mismatched types FileQueryInfoFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7276:29: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type IOStreamSpliceFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7281:10: invalid operation: i == 0 (mismatched types IOStreamSpliceFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7288:11: invalid operation: i != 0 (mismatched types IOStreamSpliceFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7289:16: invalid operation: i - 1 (mismatched types IOStreamSpliceFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7325:25: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type MountMountFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7330:10: invalid operation: m == 0 (mismatched types MountMountFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7337:11: invalid operation: m != 0 (mismatched types MountMountFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7338:16: invalid operation: m - 1 (mismatched types MountMountFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7371:27: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type MountUnmountFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7376:10: invalid operation: m == 0 (mismatched types MountUnmountFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7383:11: invalid operation: m != 0 (mismatched types MountUnmountFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7384:16: invalid operation: m - 1 (mismatched types MountUnmountFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7420:33: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type OutputStreamSpliceFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7425:10: invalid operation: o == 0 (mismatched types OutputStreamSpliceFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7432:11: invalid operation: o != 0 (mismatched types OutputStreamSpliceFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7433:16: invalid operation: o - 1 (mismatched types OutputStreamSpliceFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7472:33: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type ResolverNameLookupFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7477:10: invalid operation: r == 0 (mismatched types ResolverNameLookupFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7484:11: invalid operation: r != 0 (mismatched types ResolverNameLookupFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7485:16: invalid operation: r - 1 (mismatched types ResolverNameLookupFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7522:23: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type ResourceFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7527:10: invalid operation: r == 0 (mismatched types ResourceFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7534:11: invalid operation: r != 0 (mismatched types ResourceFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7535:16: invalid operation: r - 1 (mismatched types ResourceFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7567:29: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type ResourceLookupFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7572:10: invalid operation: r == 0 (mismatched types ResourceLookupFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7579:11: invalid operation: r != 0 (mismatched types ResourceLookupFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7580:16: invalid operation: r - 1 (mismatched types ResourceLookupFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7631:27: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type SettingsBindFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7636:10: invalid operation: s == 0 (mismatched types SettingsBindFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7643:11: invalid operation: s != 0 (mismatched types SettingsBindFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7644:16: invalid operation: s - 1 (mismatched types SettingsBindFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7696:24: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type SocketMsgFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7701:10: invalid operation: s == 0 (mismatched types SocketMsgFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7708:11: invalid operation: s != 0 (mismatched types SocketMsgFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7709:16: invalid operation: s - 1 (mismatched types SocketMsgFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7783:25: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type SubprocessFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7788:10: invalid operation: s == 0 (mismatched types SubprocessFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7795:11: invalid operation: s != 0 (mismatched types SubprocessFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7796:16: invalid operation: s - 1 (mismatched types SubprocessFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7844:23: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type TestDBusFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7849:10: invalid operation: t == 0 (mismatched types TestDBusFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7856:11: invalid operation: t != 0 (mismatched types TestDBusFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7857:16: invalid operation: t - 1 (mismatched types TestDBusFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7916:29: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type TLSCertificateFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7921:10: invalid operation: t == 0 (mismatched types TLSCertificateFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7928:11: invalid operation: t != 0 (mismatched types TLSCertificateFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7929:16: invalid operation: t - 1 (mismatched types TLSCertificateFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7975:32: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type TLSDatabaseVerifyFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7980:10: invalid operation: t == 0 (mismatched types TLSDatabaseVerifyFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7987:11: invalid operation: t != 0 (mismatched types TLSDatabaseVerifyFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:7988:16: invalid operation: t - 1 (mismatched types TLSDatabaseVerifyFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:8035:26: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type TLSPasswordFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:8040:10: invalid operation: t == 0 (mismatched types TLSPasswordFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:8047:11: invalid operation: t != 0 (mismatched types TLSPasswordFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gio/v2/gio.go:8048:16: invalid operation: t - 1 (mismatched types TLSPasswordFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/antialias.go:6:8: could not import C (no metadata for C)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:58:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:59:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:60:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:61:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:62:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:63:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:64:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:65:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:66:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:67:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:68:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:69:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:70:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:71:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:72:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:73:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:74:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:75:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:76:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:77:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:78:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:79:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:80:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:81:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:82:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:83:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:84:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:85:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:86:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:87:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:88:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:89:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:90:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:91:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/cairo/status.go:92:2: duplicate key unknown in map literal
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdkpixbuf/v2/gdkpixbuf.go:53:8: could not import C (no metadata for C)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdkpixbuf/v2/gdkpixbuf.go:117:20: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type Colorspace
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdkpixbuf/v2/gdkpixbuf.go:168:20: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type InterpType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdkpixbuf/v2/gdkpixbuf.go:216:25: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type PixbufAlphaMode
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdkpixbuf/v2/gdkpixbuf.go:256:21: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type PixbufError
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdkpixbuf/v2/gdkpixbuf.go:311:24: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type PixbufRotation
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdkpixbuf/v2/gdkpixbuf.go:346:10: invalid operation: p == 0 (mismatched types PixbufFormatFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdkpixbuf/v2/gdkpixbuf.go:353:11: invalid operation: p != 0 (mismatched types PixbufFormatFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdkpixbuf/v2/gdkpixbuf.go:354:16: invalid operation: p - 1 (mismatched types PixbufFormatFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:174:8: could not import C (no metadata for C)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:378:19: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type Alignment
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:491:18: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type AttrType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:656:23: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type BaselineShift
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:732:18: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type BidiType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:843:23: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type CoverageLevel
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:902:19: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type Direction
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:946:23: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type EllipsizeMode
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:982:19: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type FontScale
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:1034:17: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type Gravity
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:1229:21: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type GravityHint
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:1260:32: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type LayoutDeserializeError
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:1302:18: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type Overline
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:1335:20: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type RenderPart
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:1609:16: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type Script
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:1976:17: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type Stretch
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:2018:15: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type Style
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:2055:18: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type TabAlign
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:2091:23: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type TextTransform
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:2148:19: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type Underline
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:2204:17: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type Variant
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:2263:16: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type Weight
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:2318:18: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type WrapMode
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:2359:18: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type FontMask
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:2364:10: invalid operation: f == 0 (mismatched types FontMask and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:2371:11: invalid operation: f != 0 (mismatched types FontMask and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:2372:16: invalid operation: f - 1 (mismatched types FontMask and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:2422:32: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type LayoutDeserializeFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:2427:10: invalid operation: l == 0 (mismatched types LayoutDeserializeFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:2434:11: invalid operation: l != 0 (mismatched types LayoutDeserializeFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:2435:16: invalid operation: l - 1 (mismatched types LayoutDeserializeFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:2474:30: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type LayoutSerializeFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:2479:10: invalid operation: l == 0 (mismatched types LayoutSerializeFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:2486:11: invalid operation: l != 0 (mismatched types LayoutSerializeFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:2487:16: invalid operation: l - 1 (mismatched types LayoutSerializeFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:2527:20: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type ShapeFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:2532:10: invalid operation: s == 0 (mismatched types ShapeFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:2539:11: invalid operation: s != 0 (mismatched types ShapeFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:2540:16: invalid operation: s - 1 (mismatched types ShapeFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:2579:19: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type ShowFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:2584:10: invalid operation: s == 0 (mismatched types ShowFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:2591:11: invalid operation: s != 0 (mismatched types ShowFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/pango/pango.go:2592:16: invalid operation: s - 1 (mismatched types ShowFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:114:8: could not import C (no metadata for C)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:2652:17: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type AxisUse
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:2719:22: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type CrossingMode
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:2761:26: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type DevicePadFeature
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:2802:24: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type DeviceToolType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:2844:21: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type DmabufError
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:2886:26: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type DragCancelReason
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:2977:19: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type EventType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:3060:24: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type FullscreenMode
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:3092:17: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type GLError
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:3154:17: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type Gravity
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:3210:21: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type InputSource
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:3251:18: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type KeyMatch
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:3367:22: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type MemoryFormat
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:3474:20: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type NotifyType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:3515:25: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type ScrollDirection
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:3563:20: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type ScrollUnit
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:3598:24: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type SubpixelLayout
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:3644:21: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type SurfaceEdge
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:3688:22: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type TextureError
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:3728:25: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type TitlebarGesture
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:3778:30: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type TouchpadGesturePhase
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:3809:21: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type VulkanError
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:3876:21: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type AnchorHints
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:3881:10: invalid operation: a == 0 (mismatched types AnchorHints and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:3888:11: invalid operation: a != 0 (mismatched types AnchorHints and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:3889:16: invalid operation: a - 1 (mismatched types AnchorHints and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:3955:19: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type AxisFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:3960:10: invalid operation: a == 0 (mismatched types AxisFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:3967:11: invalid operation: a != 0 (mismatched types AxisFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:3968:16: invalid operation: a - 1 (mismatched types AxisFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:4028:20: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type DragAction
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:4033:10: invalid operation: d == 0 (mismatched types DragAction and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:4040:11: invalid operation: d != 0 (mismatched types DragAction and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:4041:16: invalid operation: d - 1 (mismatched types DragAction and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:4129:25: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type FrameClockPhase
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:4134:10: invalid operation: f == 0 (mismatched types FrameClockPhase and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:4141:11: invalid operation: f != 0 (mismatched types FrameClockPhase and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:4142:16: invalid operation: f - 1 (mismatched types FrameClockPhase and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:4188:15: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type GLAPI
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:4193:10: invalid operation: g == 0 (mismatched types GLAPI and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:4200:11: invalid operation: g != 0 (mismatched types GLAPI and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:4201:16: invalid operation: g - 1 (mismatched types GLAPI and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:4268:22: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type ModifierType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:4273:10: invalid operation: m == 0 (mismatched types ModifierType and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:4280:11: invalid operation: m != 0 (mismatched types ModifierType and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:4281:16: invalid operation: m - 1 (mismatched types ModifierType and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:4341:24: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type PaintableFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:4346:10: invalid operation: p == 0 (mismatched types PaintableFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:4353:11: invalid operation: p != 0 (mismatched types PaintableFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:4354:16: invalid operation: p - 1 (mismatched types PaintableFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:4400:26: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type SeatCapabilities
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:4405:10: invalid operation: s == 0 (mismatched types SeatCapabilities and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:4412:11: invalid operation: s != 0 (mismatched types SeatCapabilities and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:4413:16: invalid operation: s - 1 (mismatched types SeatCapabilities and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:4496:23: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type ToplevelState
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:4501:10: invalid operation: t == 0 (mismatched types ToplevelState and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:4508:11: invalid operation: t != 0 (mismatched types ToplevelState and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gdk/v4/gdk.go:4509:16: invalid operation: t - 1 (mismatched types ToplevelState and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/graphene/graphene.go:20:8: could not import C (no metadata for C)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gsk/v4/gsk.go:30:8: could not import C (no metadata for C)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gsk/v4/gsk.go:211:19: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type BlendMode
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gsk/v4/gsk.go:269:16: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type Corner
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gsk/v4/gsk.go:316:18: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type FillRule
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gsk/v4/gsk.go:358:23: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type GLUniformType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gsk/v4/gsk.go:410:17: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type LineCap
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gsk/v4/gsk.go:451:18: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type LineJoin
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gsk/v4/gsk.go:485:18: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type MaskMode
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gsk/v4/gsk.go:530:23: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type PathDirection
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gsk/v4/gsk.go:579:23: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type PathOperation
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gsk/v4/gsk.go:675:24: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type RenderNodeType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gsk/v4/gsk.go:765:23: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type ScalingFilter
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gsk/v4/gsk.go:796:28: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type SerializationError
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gsk/v4/gsk.go:861:27: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type TransformCategory
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gsk/v4/gsk.go:906:26: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type PathForEachFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gsk/v4/gsk.go:911:10: invalid operation: p == 0 (mismatched types PathForEachFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gsk/v4/gsk.go:918:11: invalid operation: p != 0 (mismatched types PathForEachFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gsk/v4/gsk.go:919:16: invalid operation: p - 1 (mismatched types PathForEachFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:1522:8: could not import C (no metadata for C)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:2758:40: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type AccessibleAnnouncementPriority
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:2799:32: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type AccessibleAutocomplete
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:2837:32: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type AccessibleInvalidState
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:2870:33: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type AccessiblePlatformState
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:2959:28: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type AccessibleProperty
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:3096:28: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type AccessibleRelation
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:3374:24: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type AccessibleRole
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:3574:24: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type AccessibleSort
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:3629:25: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type AccessibleState
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:3687:37: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type AccessibleTextContentChange
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:3729:35: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type AccessibleTextGranularity
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:3767:28: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type AccessibleTristate
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:3821:15: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type Align
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:3861:19: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type ArrowType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:3918:27: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type AssistantPageType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:3960:26: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type BaselinePosition
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:4004:21: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type BorderStyle
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:4086:22: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type BuilderError
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:4165:21: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type ButtonsType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:4199:31: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type CellRendererAccelMode
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:4230:26: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type CellRendererMode
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:4264:19: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type Collation
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:4319:29: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type ConstraintAttribute
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:4367:28: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type ConstraintRelation
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:4404:28: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type ConstraintStrength
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:4442:34: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type ConstraintVflParserError
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:4502:20: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type ContentFit
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:4543:20: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type CornerType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:4657:20: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type DeleteType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:4701:21: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type DialogError
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:4749:23: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type DirectionType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:4800:28: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type EditableProperties
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:4842:27: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type EntryIconPosition
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:4870:28: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type EventSequenceState
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:4904:27: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type FileChooserAction
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:4939:26: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type FileChooserError
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:4997:22: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type FilterChange
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:5034:21: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type FilterMatch
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:5070:19: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type FontLevel
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:5100:32: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type GraphicsOffloadEnabled
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:5135:18: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type IconSize
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:5163:24: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type IconThemeError
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:5209:30: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type IconViewDropPosition
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:5255:19: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type ImageType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:5321:22: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type InputPurpose
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:5372:29: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type InscriptionOverflow
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:5406:23: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type Justification
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:5439:22: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type LevelBarMode
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:5501:17: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type License
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:5570:25: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type ListTabBehavior
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:5604:21: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type MessageType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:5653:22: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type MovementStep
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:5704:25: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type NaturalWrapMode
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:5732:21: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type NotebookTab
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:5771:24: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type NumberUpLayout
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:5815:18: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type Ordering
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:5845:21: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type Orientation
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:5876:18: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type Overflow
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:5904:18: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type PackType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:5932:23: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type PadActionType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:5964:25: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type PageOrientation
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:5996:17: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type PageSet
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:6028:22: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type PanDirection
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:6068:20: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type PolicyType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:6105:22: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type PositionType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:6137:21: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type PrintDuplex
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:6171:20: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type PrintError
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:6227:30: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type PrintOperationAction
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:6265:30: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type PrintOperationResult
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:6299:20: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type PrintPages
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:6333:22: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type PrintQuality
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:6383:21: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type PrintStatus
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:6426:26: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type PropagationLimit
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:6464:26: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type PropagationPhase
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:6508:28: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type RecentManagerError
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:6579:22: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type ResponseType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:6640:32: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type RevealerTransitionType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:6690:20: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type ScrollStep
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:6752:20: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type ScrollType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:6807:26: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type ScrollablePolicy
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:6844:23: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type SelectionMode
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:6878:25: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type SensitivityType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:6911:23: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type ShortcutScope
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:6967:22: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type ShortcutType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7014:23: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type SizeGroupMode
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7048:25: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type SizeRequestMode
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7076:18: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type SortType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7112:22: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type SorterChange
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7147:21: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type SorterOrder
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7181:32: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type SpinButtonUpdatePolicy
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7218:18: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type SpinType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7310:29: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type StackTransitionType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7381:31: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type StringFilterMatchMode
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7416:23: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type SymbolicColor
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7465:23: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type SystemSetting
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7499:23: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type TextDirection
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7530:29: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type TextExtendSelection
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7558:23: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type TextViewLayer
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7592:24: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type TextWindowType
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7632:30: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type TreeViewColumnSizing
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7666:30: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type TreeViewDropPosition
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7700:27: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type TreeViewGridLines
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7734:14: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type Unit
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7771:18: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Enum() (value of type int) to type WrapMode
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7810:33: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type ApplicationInhibitFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7815:10: invalid operation: a == 0 (mismatched types ApplicationInhibitFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7822:11: invalid operation: a != 0 (mismatched types ApplicationInhibitFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7823:16: invalid operation: a - 1 (mismatched types ApplicationInhibitFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7866:29: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type BuilderClosureFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7871:10: invalid operation: b == 0 (mismatched types BuilderClosureFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7878:11: invalid operation: b != 0 (mismatched types BuilderClosureFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7879:16: invalid operation: b - 1 (mismatched types BuilderClosureFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7922:27: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type CellRendererState
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7927:10: invalid operation: c == 0 (mismatched types CellRendererState and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7934:11: invalid operation: c != 0 (mismatched types CellRendererState and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:7935:16: invalid operation: c - 1 (mismatched types CellRendererState and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8017:20: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type DebugFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8022:10: invalid operation: d == 0 (mismatched types DebugFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8029:11: invalid operation: d != 0 (mismatched types DebugFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8030:16: invalid operation: d - 1 (mismatched types DebugFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8101:21: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type DialogFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8106:10: invalid operation: d == 0 (mismatched types DialogFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8113:11: invalid operation: d != 0 (mismatched types DialogFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8114:16: invalid operation: d - 1 (mismatched types DialogFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8160:36: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type EventControllerScrollFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8165:10: invalid operation: e == 0 (mismatched types EventControllerScrollFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8172:11: invalid operation: e != 0 (mismatched types EventControllerScrollFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8173:16: invalid operation: e - 1 (mismatched types EventControllerScrollFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8226:26: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type FontChooserLevel
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8231:10: invalid operation: f == 0 (mismatched types FontChooserLevel and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8238:11: invalid operation: f != 0 (mismatched types FontChooserLevel and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8239:16: invalid operation: f - 1 (mismatched types FontChooserLevel and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8284:25: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type IconLookupFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8289:10: invalid operation: i == 0 (mismatched types IconLookupFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8296:11: invalid operation: i != 0 (mismatched types IconLookupFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8297:16: invalid operation: i - 1 (mismatched types IconLookupFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8369:20: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type InputHints
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8374:10: invalid operation: i == 0 (mismatched types InputHints and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8381:11: invalid operation: i != 0 (mismatched types InputHints and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8382:16: invalid operation: i - 1 (mismatched types InputHints and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8441:25: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type ListScrollFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8446:10: invalid operation: l == 0 (mismatched types ListScrollFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8453:11: invalid operation: l != 0 (mismatched types ListScrollFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8454:16: invalid operation: l - 1 (mismatched types ListScrollFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8493:19: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type PickFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8498:10: invalid operation: p == 0 (mismatched types PickFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8505:11: invalid operation: p != 0 (mismatched types PickFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8506:16: invalid operation: p - 1 (mismatched types PickFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8544:26: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type PopoverMenuFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8549:10: invalid operation: p == 0 (mismatched types PopoverMenuFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8556:11: invalid operation: p != 0 (mismatched types PopoverMenuFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8557:16: invalid operation: p - 1 (mismatched types PopoverMenuFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8593:29: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type ShortcutActionFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8598:10: invalid operation: s == 0 (mismatched types ShortcutActionFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8605:11: invalid operation: s != 0 (mismatched types ShortcutActionFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8606:16: invalid operation: s - 1 (mismatched types ShortcutActionFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8669:20: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type StateFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8674:10: invalid operation: s == 0 (mismatched types StateFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8681:11: invalid operation: s != 0 (mismatched types StateFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8682:16: invalid operation: s - 1 (mismatched types StateFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8754:32: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type StyleContextPrintFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8759:10: invalid operation: s == 0 (mismatched types StyleContextPrintFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8766:11: invalid operation: s != 0 (mismatched types StyleContextPrintFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8767:16: invalid operation: s - 1 (mismatched types StyleContextPrintFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8814:25: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type TextSearchFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8819:10: invalid operation: t == 0 (mismatched types TextSearchFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8826:11: invalid operation: t != 0 (mismatched types TextSearchFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8827:16: invalid operation: t - 1 (mismatched types TextSearchFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8869:24: cannot convert coreglib.ValueFromNative(unsafe.Pointer(p)).Flags() (value of type uint) to type TreeModelFlags
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8874:10: invalid operation: t == 0 (mismatched types TreeModelFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8881:11: invalid operation: t != 0 (mismatched types TreeModelFlags and untyped int)
/home/ubuntu/go/pkg/mod/github.com/diamondburned/gotk4/pkg@v0.3.1/gtk/v4/gtk.go:8882:16: invalid operation: t - 1 (mismatched types TreeModelFlags and untyped int)
/tmp/ksec-src-3187968579/download.go:73:22: cannot use gtk.OrientationVertical (constant unknown with invalid type) as gtk.Orientation value in argument to gtk.NewBox

For details on package patterns, see https://pkg.go.dev/cmd/go#hdr-Package_lists_and_patterns.

## 📋 Open PRs

**kairos-io/cluster-api-provider-kairos**

- [#38 Bump golang.org/x/oauth2 from 0.21.0 to 0.27.0 in the go_modules group across 1 directory](https://github.com/kairos-io/cluster-api-provider-kairos/pull/38) — human
**kairos-io/entangle**

- [#10 Bump the go_modules group across 1 directory with 7 updates](https://github.com/kairos-io/entangle/pull/10) — human
**kairos-io/entangle-proxy**

- [#14 Bump the go_modules group across 1 directory with 7 updates](https://github.com/kairos-io/entangle-proxy/pull/14) — human
**kairos-io/go-nodepair**

- [#27 Bump the go_modules group across 1 directory with 6 updates](https://github.com/kairos-io/go-nodepair/pull/27) — human
**kairos-io/kcrypt**

- [#509 Bump github.com/docker/docker from 27.5.1+incompatible to 28.0.0+incompatible in the go_modules group across 1 directory](https://github.com/kairos-io/kcrypt/pull/509) — human
**kairos-io/netboot**

- [#36 Bump golang.org/x/crypto from 0.39.0 to 0.45.0](https://github.com/kairos-io/netboot/pull/36) — human
**kairos-io/simple-mdns-server**

- [#4 Bump the go_modules group across 1 directory with 2 updates](https://github.com/kairos-io/simple-mdns-server/pull/4) — human
**kairos-io/tpm-helpers**

- [#6 Bump the go_modules group across 1 directory with 3 updates](https://github.com/kairos-io/tpm-helpers/pull/6) — human
**mudler/edgevpn**

- [#905 chore(deps): bump dependabot/fetch-metadata from 2.3.0 to 2.4.0](https://github.com/mudler/edgevpn/pull/905) — human
- [#923 chore(deps): bump github.com/miekg/dns from 1.1.64 to 1.1.68](https://github.com/mudler/edgevpn/pull/923) — human
- [#927 chore(deps): bump actions/checkout from 4 to 5](https://github.com/mudler/edgevpn/pull/927) — human
- [#939 chore(deps): bump actions/setup-go from 5 to 6](https://github.com/mudler/edgevpn/pull/939) — human
- [#942 chore(deps): bump github.com/onsi/gomega from 1.37.0 to 1.38.2](https://github.com/mudler/edgevpn/pull/942) — human
- [#943 chore(deps): bump codecov/codecov-action from 5.5.0 to 5.5.1](https://github.com/mudler/edgevpn/pull/943) — human
- [#946 chore(deps): bump github.com/libp2p/go-libp2p-kad-dht from 0.33.1 to 0.35.1](https://github.com/mudler/edgevpn/pull/946) — human
- [#951 chore(deps): bump actions/download-artifact from 5 to 6](https://github.com/mudler/edgevpn/pull/951) — human
- [#961 chore(deps): bump github.com/libp2p/go-libp2p-pubsub from 0.14.2 to 0.15.0](https://github.com/mudler/edgevpn/pull/961) — human
- [#1006 chore(deps): bump github.com/labstack/echo/v4 from 4.13.3 to 4.15.1](https://github.com/mudler/edgevpn/pull/1006) — human
- [#1009 chore(deps): bump docs/themes/docsy from `bbf68d4` to `01c827e`](https://github.com/mudler/edgevpn/pull/1009) — human
**mudler/entities**

- [#10 Bump golang.org/x/text from 0.3.2 to 0.3.8](https://github.com/mudler/entities/pull/10) — human
- [#11 Bump golang.org/x/net from 0.0.0-20191209160850-c0dbc17a3553 to 0.7.0](https://github.com/mudler/entities/pull/11) — human
- [#12 Bump golang.org/x/sys from 0.0.0-20200102141924-c96a22e43c9c to 0.1.0](https://github.com/mudler/entities/pull/12) — human

## 🤖 Bot PR ledger

_No bot PRs yet._

