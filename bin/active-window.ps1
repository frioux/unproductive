$job = Start-Job -ScriptBlock {

    Add-Type @"
      using System;
      using System.Runtime.InteropServices;
      using System.Text;
      public class Tricks {
        [DllImport("user32.dll")]
        public static extern IntPtr GetForegroundWindow();
        [DllImport("user32.dll")]
        public static extern int GetWindowText(IntPtr hWnd, StringBuilder text, int count);
        }
"@

    function GetActiveWindow {
        # Set-Variable nChars -option ReadOnly -value 256
        $stringBuilder = New-Object System.Text.StringBuilder(256)
        $handle = [tricks]::GetForegroundWindow()
    
        if ([tricks]::GetWindowText($handle, $stringBuilder, 256) -gt 0) {     
            return $stringBuilder.ToString()    
        }
        return $null
    }

    function GetSSID {
        $output = netsh.exe wlan show interfaces
        $ssidsearchresult = $output | Select-String -Pattern 'SSID'
        if ($ssidsearchresult.length -lt 1) {
            return $null
        }
        else
        {
            return ($ssidsearchresult -split ":")[1].Trim()
        }    
    }

    function IsConnectedToOpenVPN{
        $result = Get-Process -Name "openvpn" -ErrorAction SilentlyContinue
        return $result -ne $null
    }

    $final = [pscustomobject] @{
        Time = [int][double]::Parse((Get-Date -UFormat %s))
        SSID = GetSSID
        VPN = IsConnectedToOpenVPN
        ActiveWindow = GetActiveWindow

    }

    Write-Output ([string]::join("`t", ($final.Time, $final.SSID, $final.VPN, $final.ActiveWindow)))
}

Out-Null -inputObject (Wait-Job $job)
Receive-Job $job