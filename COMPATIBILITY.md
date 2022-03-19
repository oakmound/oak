# Oak Compatibility Matrix

| OS / Target  | Rendering   | Multi-Window | Audio       | Joysticks         | Cgo Required | Architectures Supported |
|:------------:|:-----------:|:------------:|:-----------:|:-----------------:|:------------:|:-----------------------:|
| windows      | Yes         | Yes          | Yes         | XInput            | No           | 386, amd64, arm64       |
| linux        | Yes         | Yes          | Yes         | XBox 360 (#133)   | No           | amd64, arm, arm64       |
| osx (darwin) | Yes         | Yes          | Yes         | No (#87)          | Yes (#175)   | amd64                   |
| js           | Yes         | N/A          | No (#174)   | Standard Mapping  | No           | wasm                    |
| android      | Experimental| N/A          | No          | No                | Yes          | arm64                   |
| ios          | No (#49)    | N/A          | N/A         | N/A               | N/A          |                         |

## Window Options

| OS / Target  | Get Cursor Position* | Fullscreen | Borderless | Set Title** | Reposition | Window On Top | Hide Cursor | Show Notification | Set Tray Icon |
|:------------:|:--------------------:|:----------:|:----------:|:-----------:|:----------:|:-------------:|:-----------:|:-----------------:|:-------------:|
| windows      | Yes                  | Yes        | Yes        | Yes         | Yes        | Yes           | Yes         | Yes               | Yes           |
| linux        | Yes                  | Yes        | Yes        | No          | Yes        | No            | No          | No                | No            |
| osx (darwin) | Yes                  | Yes        | Yes        | No          | Yes        | No            | Yes         | No                | No            |
| wasm+js      | No                   | No         | N/A        | N/A         | N/A        | N/A           | No          | No                | No            |
| android      | No                   | Required   | Required   | No          | N/A        | N/A           | N/A         | No                | No            |

\* This refers to asking the OS where the cursor is, which can inform the absolute position of the cursor even if it is outside of the Oak window. Oak can always tell you where the cursor is if it is within the Oak window.

\*\* Changing the title of the window after it is created.

## Other Compatibility Issues

* Issue #171: Under an unknown condition, Oak fails to render or intialize on OSX.
