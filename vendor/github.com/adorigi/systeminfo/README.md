# Systeminfo
Get useful information about the system you are running, right from your terminal

## Usage

- Initialize module

    ```
    go mod init <module-name>
    ```

- Add the package

    ```
    go get github.com/adorigi/systeminfo
    ```

- Import the package to your go file

    ```
    import (
        ...
        "github.com/adorigi/systeminfo"
        ...
    )
    ```

- Collect the Stats and print them

    ```
    func main() {
        info := systeminfo.CollectStats()

        fmt.Println("CPU name: " + info.CpuName)
        fmt.Println("CPU Architecture: ", info.CpuArch)
        fmt.Println("Operating System", info.OperatingSystem)
        fmt.Println("Available Disk Storage: ", info.DiskAvailable, info.StorageUnit)
        fmt.Println("Used Disk Storage: ", info.DiskUsed, info.StorageUnit)
        fmt.Println("Disk Used %: ", info.DiskUsedPercent, "%")
        fmt.Println("Local IPv4: ", info.LocalIPv4)
        fmt.Println("Global IP: ", info.GlobalIP)

    }
    ```


### Note: Look into the example folder for clarity