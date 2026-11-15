func Commands() (commands []structs.Command) {
        // TODO Add the following commands: sharpgen
        commands = append(
                commands, cd(), createProcess(), donutCmd(), download(), env(), executeAssembly(), executePE(), executeShellcode(),
                exit(), ifconfig(), invokeAssembly(), ja3(), killdate(), killProcess(), listAssembly(), loadAssembly(), ls(),
                makeToken(), maxRetry(), memfd(), memory(), mimikatz(), netstat(), nslookup(), parrot(), pipes(), ps(), pwd(),
                rev2Self(), rm(), run(), runas(), sdelete(), shell(), skew(), sleep(), socks(), srdiCmd(), ssh(), stealToken(), socketSpy(), timeStomp(), token(), upload(), uptime(),
