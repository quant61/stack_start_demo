
When process is created
kernel creates memory mappings including code, stack, head and so on.
It also places some info there.

This demo shows how is stack memory filled on process start.

Usually there are lots of stages before main, but you can get rid of them and set entrypoint directly to your code.

But even in this case there are some work made by kernel.

This demo builds minimally possible executable, runs it in debug more and checks its memory contents.


This currently supports only linux on x86_64.

Target executable(ELF) consists only from elf header, single segment/prog header and single int0x3(debug trap) instruction.

Binary is run with debugging and demo shows what its memory contents is.

windows version is in progress(tested in wine):
- read memory(not tested)
- registers(cannot get info yet, probably bad args)
- maps()


TODO:
- more flexible, configurable and more library-like code
- running in both usual and testing mode
- not just print mappings, make them passable to other code
- more platforms(both OS and cpu)
- PE(windows) and other binary types support
- support for foreign platforms(like other cpu arch in qemu, wine or wsl)
- make use of github actions
  - generate output right inside github runners 

