package main

import (
	"encoding/base64"
	"sort"
)

const (
	file6173736574732f626173652e637373                                           = "OnJvb3QgewogIC0tYm9yZGVyLXJhZGl1czogNXB4OwogIC0tYm9yZGVyOiAxcHggc29saWQgYmxhY2s7CiAgLyogY29sb3Igc2NoZW1lIHdhcyBkZXNpZ25lZCB1c2luZyBjb2xvcm1pbmQuaW8gKi8KICAtLXByaW1hcnk6ICMyMTdjOGY7CiAgLS1zZWNvbmRhcnk6ICNjMjkyNjA7CiAgLS1oaWdobGlnaHQ6ICNhZTQ5NGE7CiAgLS1iYWNrZ3JvdW5kLWNvbG9yOiB3aGl0ZTsKICAtLWRpbW1lZDogI2M2ZGFiZTsKICAtLWxpZ2h0LWJnOiAjZjBmMGYwOwogIC0tZm9udC1jb2xvcjogYmxhY2s7Cn0KCkBtZWRpYSAocHJlZmVycy1jb2xvci1zY2hlbWU6IGRhcmspIHsKICA6cm9vdCB7CiAgICAtLWJvcmRlcjogMXB4IHNvbGlkIHdoaXRlOwogICAgLyogY29sb3Igc2NoZW1lIHdhcyBkZXNpZ25lZCB1c2luZyBjb2xvcm1pbmQuaW8gKi8KICAgIC0tcHJpbWFyeTogIzIxN2M4ZjsKICAgIC0tc2Vjb25kYXJ5OiAjYzI5MjYwOwogICAgLS1oaWdobGlnaHQ6ICNhZTQ5NGE7CiAgICAtLWJhY2tncm91bmQtY29sb3I6ICMyMTIxMjE7CiAgICAtLWRpbW1lZDogIzMzMzczMTsKICAgIC0tbGlnaHQtYmc6ICMzMDMwMzA7CiAgICAtLWZvbnQtY29sb3I6IHdoaXRlOwogIH0KfQoKYm9keSB7CiAgZm9udC1mYW1pbHk6IENhbnRhcmVsbCwgc2FuczsKICBsaW5lLWhlaWdodDogMS42ZW07CiAgbWFyZ2luOiAwOwogIGNvbG9yOiB2YXIoLS1mb250LWNvbG9yKTsKICBiYWNrZ3JvdW5kLWNvbG9yOiB2YXIoLS1iYWNrZ3JvdW5kLWNvbG9yKTsKfQoKaGVhZGVyLApmb290ZXIgewogIHBhZGRpbmc6IDJlbTsKICBiYWNrZ3JvdW5kLWNvbG9yOiB2YXIoLS1saWdodC1iZyk7Cn0KCm5hdiA+IGgxIHsKICBmb250LXdlaWdodDogYm9sZDsKICBkaXNwbGF5OiBpbmxpbmU7CiAgbWFyZ2luLXJpZ2h0OiAyZW07CiAgdGV4dC10cmFuc2Zvcm06IHVwcGVyY2FzZTsKfQoKbmF2ID4gdWwgewogIGxpc3Qtc3R5bGU6IG5vbmU7CiAgZGlzcGxheTogaW5saW5lLWZsZXg7CiAgcGFkZGluZzogMDsKICBtYXJnaW46IDA7Cn0KCm5hdiA+IHVsID4gbGkgPiBhIHsKICBmb250LXdlaWdodDogYm9sZDsKfQoKbmF2ID4gZm9ybSB7CiAgZGlzcGxheTogaW5saW5lOwogIGZsb2F0OiByaWdodDsKfQoKaW5wdXQjc2VhcmNoLXBhdHRlcm4gewogIGJvcmRlci1yYWRpdXM6IHZhcigtLWJvcmRlci1yYWRpdXMpOwogIGJvcmRlcjogdmFyKC0tYm9yZGVyKTsKICBwYWRkaW5nOiAwLjVlbSAyZW07CiAgY29sb3I6IHZhcigtLWZvbnQtY29sb3IpOwogIGJhY2tncm91bmQtY29sb3I6IHZhcigtLWxpZ2h0LWJnKTsKfQoKdGV4dGFyZWEjbm90ZSB7CiAgYm9yZGVyLXJhZGl1czogdmFyKC0tYm9yZGVyLXJhZGl1cyk7CiAgYm9yZGVyOiB2YXIoLS1ib3JkZXIpOwogIHBhZGRpbmc6IDJlbTsKICBmb250LWZhbWlseTogbW9ub3NwYWNlOwogIGZvbnQtc2l6ZTogbGFyZ2VyOwogIHdpZHRoOiAxMDAlOwogIGhlaWdodDogMTJjaDsKICB0cmFuc2l0aW9uOiBlYXNlLW91dCAwLjNzOwogIGNvbG9yOiB2YXIoLS1mb250LWNvbG9yKTsKICBiYWNrZ3JvdW5kLWNvbG9yOiB2YXIoLS1saWdodC1iZyk7Cn0KCnRleHRhcmVhI25vdGU6Zm9jdXMgewogIGJveC1zaGFkb3c6IDAgMCAwIDEwMHZ3ICMwMDAwMDBjYzsKICBoZWlnaHQ6IDMwY2g7CiAgdHJhbnNpdGlvbjogZWFzZS1pbiAwLjNzOwp9CgpidXR0b24gewogIGJvcmRlci1yYWRpdXM6IHZhcigtLWJvcmRlci1yYWRpdXMpOwogIGJvcmRlcjogdmFyKC0tYm9yZGVyKTsKICBib3JkZXItY29sb3I6IHZhcigtLXByaW1hcnkpOwogIGJhY2tncm91bmQtY29sb3I6IHZhcigtLXByaW1hcnkpOwogIGNvbG9yOiB3aGl0ZTsKICBwYWRkaW5nOiAwLjVlbSAyZW07CiAgdGV4dC10cmFuc2Zvcm06IHVwcGVyY2FzZTsKICB0cmFuc2l0aW9uOiBlYXNlLW91dCAwLjNzOwp9CgpidXR0b246aG92ZXIgewogIGJhY2tncm91bmQtY29sb3I6IHZhcigtLWhpZ2hsaWdodCk7CiAgdHJhbnNpdGlvbjogZWFzZS1pbiAwLjNzOwp9Cgpmb3JtLm5vdGUtZWRpdCA+IGJ1dHRvblt0eXBlPSJzdWJtaXQiXSB7CiAgbWFyZ2luOiAxZW0gMDsKfQoKYSB7CiAgdGV4dC1kZWNvcmF0aW9uOiBub25lOwogIHdvcmQtYnJlYWs6IGJyZWFrLXdvcmQ7CiAgY29sb3I6IHZhcigtLXByaW1hcnkpOwp9CgphOmhvdmVyIHsKICBjb2xvcjogdmFyKC0tc2Vjb25kYXJ5KTsKICB0cmFuc2l0aW9uOiBlYXNlLWluIDAuM3M7Cn0KCm1haW4gewogIHdpZHRoOiBtaW4oNzJjaCwgMTAwJSk7CiAgbWFyZ2luOiA0ZW0gYXV0bzsKfQoKbWFpbiA+IGgyIHsKICB0ZXh0LWFsaWduOiBjZW50ZXI7CiAgbWFyZ2luLXRvcDogMmVtOwp9Cgpmb3JtLm5vdGUtZWRpdCB7CiAgbWFyZ2luOiAyZW0gMWVtIDNlbTsKfQoKYXJ0aWNsZSB7CiAgbWFyZ2luOiAyZW0gMDsKfQoKYXJ0aWNsZSA+IGRpdi5jb250ZW50IHsKICBtYXJnaW4tbGVmdDogMmVtOwp9CgpkaXYuYXJ0aWNsZS1oZWFkbGluZSA+IGgzIHsKICBkaXNwbGF5OiBpbmxpbmU7Cn0KCnNwYW4ubGFzdC11cGRhdGUgewogIGNvbG9yOiB2YXIoLS1kaW1tZWQpOwogIHRyYW5zaXRpb246IGVhc2UtaW4gMC4zczsKICBmbG9hdDogcmlnaHQ7Cn0KCnNwYW4ubGFzdC11cGRhdGU6aG92ZXIgewogIGNvbG9yOiB2YXIoLS1mb250LWNvbG9yKTsKICB0cmFuc2l0aW9uOiBlYXNlLWluIDAuM3M7Cn0KCnZhciB7CiAgZm9udC1mYW1pbHk6IG1vbm9zcGFjZTsKfQoKZm9vdGVyID4gdWwgewogIGxpc3Qtc3R5bGU6IG5vbmU7CiAgZGlzcGxheTogaW5saW5lLWZsZXg7CiAganVzdGlmeS1jb250ZW50OiBzcGFjZS1hcm91bmQ7CiAgcGFkZGluZzogMDsKICBtYXJnaW46IDA7CiAgd2lkdGg6IDEwMCU7Cn0K"
	file6173736574732f66617669636f6e2e737667                                     = "PHN2ZyB3aWR0aD0iMWVtIiBoZWlnaHQ9IjFlbSIgdmlld0JveD0iMCAwIDE2IDE2IiBjbGFzcz0iYmkgYmktY2xpcGJvYXJkLWNoZWNrIiBmaWxsPSJjdXJyZW50Q29sb3IiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CiAgPHBhdGggZmlsbC1ydWxlPSJldmVub2RkIiBkPSJNNCAxLjVIM2EyIDIgMCAwIDAtMiAyVjE0YTIgMiAwIDAgMCAyIDJoMTBhMiAyIDAgMCAwIDItMlYzLjVhMiAyIDAgMCAwLTItMmgtMXYxaDFhMSAxIDAgMCAxIDEgMVYxNGExIDEgMCAwIDEtMSAxSDNhMSAxIDAgMCAxLTEtMVYzLjVhMSAxIDAgMCAxIDEtMWgxdi0xeiIvPgogIDxwYXRoIGZpbGwtcnVsZT0iZXZlbm9kZCIgZD0iTTkuNSAxaC0zYS41LjUgMCAwIDAtLjUuNXYxYS41LjUgMCAwIDAgLjUuNWgzYS41LjUgMCAwIDAgLjUtLjV2LTFhLjUuNSAwIDAgMC0uNS0uNXptLTMtMUExLjUgMS41IDAgMCAwIDUgMS41djFBMS41IDEuNSAwIDAgMCA2LjUgNGgzQTEuNSAxLjUgMCAwIDAgMTEgMi41di0xQTEuNSAxLjUgMCAwIDAgOS41IDBoLTN6bTQuMzU0IDcuMTQ2YS41LjUgMCAwIDEgMCAuNzA4bC0zIDNhLjUuNSAwIDAgMS0uNzA4IDBsLTEuNS0xLjVhLjUuNSAwIDEgMSAuNzA4LS43MDhMNy41IDkuNzkzbDIuNjQ2LTIuNjQ3YS41LjUgMCAwIDEgLjcwOCAweiIvPgo8L3N2Zz4K"
	file6d6967726174696f6e732f30315f696e697469616c5f73657475702e646f776e2e73716c = "RFJPUCBUQUJMRSBub3RlX2Z0czsKRFJPUCBUQUJMRSBub3RlOw"
	file6d6967726174696f6e732f30315f696e697469616c5f73657475702e75702e73716c     = "Q1JFQVRFIFRBQkxFIG5vdGUgKAogICAgaWQgSU5URUdFUiBQUklNQVJZIEtFWSwKICAgIC0tIFJGQzMzMzkgZm9ybWF0dGVkCiAgICBkYXRlX2NyZWF0ZWQgVEVYVCBOT1QgTlVMTCwKICAgIGRhdGVfdXBkYXRlZCBURVhULAogICAgbWFya2Rvd24gVEVYVCBOT1QgTlVMTCwKICAgIGh0bWwgVEVYVCBOT1QgTlVMTAopOwoKQ1JFQVRFIFZJUlRVQUwgVEFCTEUgbm90ZV9mdHMgVVNJTkcgZnRzNChpZCwgbWFya2Rvd24pOwpJTlNFUlQgSU5UTyBub3RlX2Z0cyBTRUxFQ1QgaWQsIG1hcmtkb3duIEZST00gbm90ZTsKLS0gU0VMRUNUICogRlJPTSBub3RlX2Z0cyBXSEVSRSBtYXJrZG93biBNQVRDSCAnbGluayc7CgpDUkVBVEUgVFJJR0dFUiBub3RlX2Z0c19iZWZvcmVfZGVsZXRlIEJFRk9SRSBERUxFVEUgT04gbm90ZQpCRUdJTgogICAgREVMRVRFIEZST00gbm90ZV9mdHMgV0hFUkUgbm90ZV9mdHMuaWQgPSBPTEQuaWQ7CkVORDsKCkNSRUFURSBUUklHR0VSIG5vdGVfZnRzX2FmdGVyX2luc2VydCBBRlRFUiBJTlNFUlQgT04gbm90ZQpCRUdJTgogICAgSU5TRVJUIElOVE8gbm90ZV9mdHMoaWQsIG1hcmtkb3duKSBWQUxVRVMoTkVXLmlkLCBORVcubWFya2Rvd24pOwpFTkQ7CgpDUkVBVEUgVFJJR0dFUiBub3RlX2Z0c19hZnRlcl91cGRhdGUgQUZURVIgVVBEQVRFIE9OIG5vdGUKQkVHSU4KICAgIFVQREFURSBub3RlX2Z0cyBTRVQgbWFya2Rvd24gPSBORVcubWFya2Rvd24gV0hFUkUgbm90ZV9mdHMuaWQgPSBORVcuaWQ7CkVORDsK"
	file76696577732f6572726f722e676f68746d6c                                     = "e3sgZGVmaW5lICJjb250ZW50IiB9fQo8cD57eyAuRXJyb3JNZXNzYWdlIH19PC9wPgp7eyBlbmQgfX0"
	file76696577732f696e6465782e676f68746d6c                                     = "e3sgZGVmaW5lICJjb250ZW50IiB9fQo8Zm9ybSBjbGFzcz0ibm90ZS1lZGl0IiBtZXRob2Q9IlBPU1QiIGFjdGlvbj0ie3sgLlN1Ym1pdEFjdGlvbiB9fSI+CiAgPHRleHRhcmVhIGlkPSJub3RlIiBuYW1lPSJub3RlIj57eyAuRWRpdFRleHQgfX08L3RleHRhcmVhPgogIDxidXR0b24gdHlwZT0ic3VibWl0Ij5TdWJtaXQ8L2J1dHRvbj4KPC9mb3JtPgp7eyByYW5nZSAkZGF5IDo9IC5EYXlzIH19CiAgPGgyPjx0aW1lIGRhdGV0aW1lPSd7eyAkZGF5LkZvcm1hdCAiMjAwNi0wMS0wMlQxNTowNDowNS4wMDAtMDcwMCIgfX0nPnt7ICRkYXkuRm9ybWF0ICJNb25kYXksIDAyLUphbi0yMDA2IiB9fTwvdGltZT48L2gyPgogIHt7IHJhbmdlICRfLCAkbm90ZSA6PSAoaW5kZXggJC5Ob3Rlc0J5RGF5ICRkYXkpIH19CiAgPGFydGljbGU+CiAgICA8ZGl2IGNsYXNzPSJhcnRpY2xlLWhlYWRsaW5lIj4KICAgICAgPGgzPjx0aW1lPnt7ICRub3RlLkRhdGVDcmVhdGVkLkZvcm1hdCAiMTU6MDQiIH19PC90aW1lPjwvaDM+IDxhIGNsYXNzPSJlZGl0IiBocmVmPSIvbm90ZS97eyAkbm90ZS5JRCB9fS9lZGl0Ij7inI/vuI88L2E+CiAgICAgIHt7IGlmIG5vdCAkbm90ZS5EYXRlVXBkYXRlZC5Jc1plcm8gfX0KICAgICAgPHNwYW4gY2xhc3M9Imxhc3QtdXBkYXRlIj5MYXN0IHVwZGF0ZSA8dGltZSBjbGFzcz0idXBkYXRlZCIgZGF0ZXRpbWU9J3t7ICRub3RlLkRhdGVVcGRhdGVkLkZvcm1hdCAiMjAwNi0wMS0wMlQxNTowNDowNS4wMDAtMDcwMCIgfX0nPnt7ICRub3RlLkRhdGVVcGRhdGVkLkZvcm1hdCAiMTU6MDQiIH19PC90aW1lPjwvc3Bhbj4KICAgICAge3sgZW5kIH19CiAgICA8L2Rpdj4KICAgIDxkaXYgY2xhc3M9ImNvbnRlbnQiPgogICAge3sgJG5vdGUuSFRNTCB9fQogICAgPC9kaXY+CiAgPC9hcnRpY2xlPgogIHt7IGVuZCB9fQp7eyBlbmQgfX0Ke3sgZW5kIH19"
	file76696577732f6c61796f7574732f626173652e676f68746d6c                       = "e3sgZGVmaW5lICJoZWFkZXIiIH19CjxuYXY+CiAgPGgxPnt7IC5BcHBOYW1lIH19PC9oMT4KICA8dWw+CiAgICA8bGk+PGEgaHJlZj0iLyI+SG9tZSAvPC9hPjwvbGk+CiAgPC91bD4KICA8Zm9ybSBjbGFzcz0ibmF2LXNlYXJjaCIgbWV0aG9kPSJHRVQiIGFjdGlvbj0iL3NlYXJjaCI+CiAgICB7ey8qIDxsYWJlbCBmb3I9InNlYXJjaC1wYXR0ZXJuIj5TZWFyY2g6PC9sYWJlbD4gKi99fQogICAgPGlucHV0IHR5cGU9InNlYXJjaCIgbWlubGVuZ3RoPTMgaWQ9InNlYXJjaC1wYXR0ZXJuIiBuYW1lPSJzZWFyY2gtcGF0dGVybiIgcGxhY2Vob2xkZXI9IldoYXQgYXJlIHlvdSBsb29raW5nIGZvcj8iPgogICAge3svKiA8YnV0dG9uIHR5cGU9InN1Ym1pdCI+U2VhcmNoPC9idXR0b24+ICovfX0KICA8L2Zvcm0+CjwvbmF2Pgp7eyBlbmQgfX0KCnt7IGRlZmluZSAiZm9vdGVyIiB9fQo8dWwgY2xhc3M9ImZvb3RlciI+CiAgPGxpPnt7IC5BcHBOYW1lIH19PC9saT4KICA8bGk+VmVyc2lvbjogPHZhcj57eyAuVmVyc2lvbiB9fTwvdmFyPjwvbGk+CiAgPGxpPnJlbmRlcmVkIGF0IDx0aW1lIGRhdGV0aW1lPSd7eyAuUmVuZGVyRGF0ZS5Gb3JtYXQgIjIwMDYtMDEtMDJUMTU6MDQ6MDUuMDAwLTA3MDAiIH19Jz57eyAuUmVuZGVyRGF0ZS5Gb3JtYXQgIjIwMDYtMDEtMDJUMTU6MDQiIH19PC90aW1lPjwvbGk+CjwvdWw+Cnt7IGVuZCB9fQoKe3sgZGVmaW5lICJtYWluIiB9fQo8aDI+e3sgLkhlYWRpbmcgfX08L2gyPgp7eyB0ZW1wbGF0ZSAiY29udGVudCIgLkNvbnRlbnQgfX0Ke3sgZW5kIH19Cgp7eyBkZWZpbmUgImJhc2UiIH19CjwhZG9jdHlwZSBodG1sPgo8aHRtbCBsYW5nPSJlbiI+CjxoZWFkPgogIDx0aXRsZT57eyAuVGl0bGUgfX08L3RpdGxlPgogIDxtZXRhIGNoYXJzZXQ9InV0Zi04Ij4KICA8bWV0YSBuYW1lPSJ2aWV3cG9ydCIgY29udGVudD0id2lkdGg9ZGV2aWNlLXdpZHRoLCBpbml0aWFsLXNjYWxlPTEiPgogIDxsaW5rIHJlbD0iaWNvbiIgaHJlZj0iL2Fzc2V0cy9mYXZpY29uLnN2ZyI+CiAgPGxpbmsgcmVsPSJzdHlsZXNoZWV0IiBocmVmPSIvL2NkbmpzLmNsb3VkZmxhcmUuY29tL2FqYXgvbGlicy9oaWdobGlnaHQuanMvMTAuMS4yL3N0eWxlcy9kZWZhdWx0Lm1pbi5jc3MiPgogIDxsaW5rIHJlbD0ic3R5bGVzaGVldCIgaHJlZj0iL2Fzc2V0cy9iYXNlLmNzcyI+CiAgPHNjcmlwdCBkZWZlciBzcmM9Ii8vY2RuanMuY2xvdWRmbGFyZS5jb20vYWpheC9saWJzL2hpZ2hsaWdodC5qcy8xMC4xLjIvaGlnaGxpZ2h0Lm1pbi5qcyI+PC9zY3JpcHQ+CjwvaGVhZD4KPGJvZHk+CiAgPGhlYWRlcj57eyB0ZW1wbGF0ZSAiaGVhZGVyIiAuSGVhZGVyIH19PC9oZWFkZXI+CiAgPG1haW4+e3sgdGVtcGxhdGUgIm1haW4iIC5NYWluIH19PC9tYWluPgogIDxmb290ZXI+e3sgdGVtcGxhdGUgImZvb3RlciIgLkZvb3RlciB9fTwvZm9vdGVyPgogIDxzY3JpcHQgc3JjPSIvL2NkbmpzLmNsb3VkZmxhcmUuY29tL2FqYXgvbGlicy9oaWdobGlnaHQuanMvMTAuMS4yL2hpZ2hsaWdodC5taW4uanMiPjwvc2NyaXB0PgogIDxzY3JpcHQ+aGxqcy5pbml0SGlnaGxpZ2h0aW5nT25Mb2FkKCk7PC9zY3JpcHQ+CjwvYm9keT4KPC9odG1sPgp7eyBlbmQgfX0KCnt7IHRlbXBsYXRlICJiYXNlIiAuIH19Cg"
)

// Embedded implements github.com/klingtnet/embed/Embed .
type Embedded struct {
	embedMap map[string]string
}

// Embeds stores the embedded data.
var Embeds = Embedded{
	embedMap: map[string]string{
		"assets/base.css":                      file6173736574732f626173652e637373,
		"assets/favicon.svg":                   file6173736574732f66617669636f6e2e737667,
		"migrations/01_initial_setup.down.sql": file6d6967726174696f6e732f30315f696e697469616c5f73657475702e646f776e2e73716c,
		"migrations/01_initial_setup.up.sql":   file6d6967726174696f6e732f30315f696e697469616c5f73657475702e75702e73716c,
		"views/error.gohtml":                   file76696577732f6572726f722e676f68746d6c,
		"views/index.gohtml":                   file76696577732f696e6465782e676f68746d6c,
		"views/layouts/base.gohtml":            file76696577732f6c61796f7574732f626173652e676f68746d6c,
	},
}

// Files implements github.com/klingtnet/embed/Embed .
func (e Embedded) Files() []string {
	var fs []string
	for f := range e.embedMap {
		fs = append(fs, f)
	}
	sort.Strings(fs)
	return fs
}

// File implements github.com/klingtnet/embed/Embed .
func (e Embedded) File(path string) []byte {
	file, ok := e.embedMap[path]
	if !ok {
		return nil
	}
	d, err := base64.RawStdEncoding.DecodeString(file)
	if err != nil {
		panic(err)
	}
	return d
}

// FileString implements github.com/klingtnet/embed/Embed .
func (e Embedded) FileString(path string) string {
	return string(e.File(path))
}
