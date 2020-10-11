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
	file76696577732f6c61796f7574732f626173652e676f68746d6c                       = "e3sgZGVmaW5lICJoZWFkZXIiIH19CjxuYXY+CiAgPGgxPnt7IC5BcHBOYW1lIH19PC9oMT4KICA8dWw+CiAgICA8bGk+PGEgaHJlZj0iLyI+SG9tZSAvPC9hPjwvbGk+CiAgPC91bD4KICA8Zm9ybSBjbGFzcz0ibmF2LXNlYXJjaCIgbWV0aG9kPSJHRVQiIGFjdGlvbj0iL3NlYXJjaCI+CiAgICA8bGFiZWwgZm9yPSJzZWFyY2gtcGF0dGVybiI+U2VhcmNoOjwvbGFiZWw+CiAgICA8aW5wdXQgdHlwZT0ic2VhcmNoIiBtaW5sZW5ndGg9MyBpZD0ic2VhcmNoLXBhdHRlcm4iIG5hbWU9InNlYXJjaC1wYXR0ZXJuIj4KICAgIHt7LyogPGJ1dHRvbiB0eXBlPSJzdWJtaXQiPlNlYXJjaDwvYnV0dG9uPiAqL319CiAgPC9mb3JtPgo8L25hdj4Ke3sgZW5kIH19Cgp7eyBkZWZpbmUgImZvb3RlciIgfX0KPHVsIGNsYXNzPSJmb290ZXIiPgogIDxsaT57eyAuQXBwTmFtZSB9fTwvbGk+CiAgPGxpPlZlcnNpb246IDx2YXI+e3sgLlZlcnNpb24gfX08L3Zhcj48L2xpPgogIDxsaT5yZW5kZXJlZCBhdCA8dGltZSBkYXRldGltZT0ne3sgLlJlbmRlckRhdGUuRm9ybWF0ICIyMDA2LTAxLTAyVDE1OjA0OjA1LjAwMC0wNzAwIiB9fSc+e3sgLlJlbmRlckRhdGUuRm9ybWF0ICIyMDA2LTAxLTAyVDE1OjA0IiB9fTwvdGltZT48L2xpPgo8L3VsPgp7eyBlbmQgfX0KCnt7IGRlZmluZSAibWFpbiIgfX0KPGgyPnt7IC5IZWFkaW5nIH19PC9oMj4Ke3sgdGVtcGxhdGUgImNvbnRlbnQiIC5Db250ZW50IH19Cnt7IGVuZCB9fQoKe3sgZGVmaW5lICJiYXNlIiB9fQo8IWRvY3R5cGUgaHRtbD4KPGh0bWwgbGFuZz0iZW4iPgo8aGVhZD4KICA8dGl0bGU+e3sgLlRpdGxlIH19PC90aXRsZT4KICA8bWV0YSBjaGFyc2V0PSJ1dGYtOCI+CiAgPG1ldGEgbmFtZT0idmlld3BvcnQiIGNvbnRlbnQ9IndpZHRoPWRldmljZS13aWR0aCwgaW5pdGlhbC1zY2FsZT0xIj4KICA8bGluayByZWw9Imljb24iIGhyZWY9Ii9hc3NldHMvZmF2aWNvbi5zdmciPgogIDxsaW5rIHJlbD0ic3R5bGVzaGVldCIgaHJlZj0iLy9jZG5qcy5jbG91ZGZsYXJlLmNvbS9hamF4L2xpYnMvaGlnaGxpZ2h0LmpzLzEwLjEuMi9zdHlsZXMvZGVmYXVsdC5taW4uY3NzIj4KICA8bGluayByZWw9InN0eWxlc2hlZXQiIGhyZWY9Ii9hc3NldHMvYmFzZS5jc3MiPgogIDxsaW5rIHJlbD0ic3R5bGVzaGVldCIgaHJlZj0iLy9jZG5qcy5jbG91ZGZsYXJlLmNvbS9hamF4L2xpYnMvaGlnaGxpZ2h0LmpzLzEwLjEuMi9zdHlsZXMvZGVmYXVsdC5taW4uY3NzIj4KICA8c2NyaXB0IGRlZmVyIHNyYz0iLy9jZG5qcy5jbG91ZGZsYXJlLmNvbS9hamF4L2xpYnMvaGlnaGxpZ2h0LmpzLzEwLjEuMi9oaWdobGlnaHQubWluLmpzIj48L3NjcmlwdD4KPC9oZWFkPgo8Ym9keT4KICA8aGVhZGVyPnt7IHRlbXBsYXRlICJoZWFkZXIiIC5IZWFkZXIgfX08L2hlYWRlcj4KICA8bWFpbj57eyB0ZW1wbGF0ZSAibWFpbiIgLk1haW4gfX08L21haW4+CiAgPGZvb3Rlcj57eyB0ZW1wbGF0ZSAiZm9vdGVyIiAuRm9vdGVyIH19PC9mb290ZXI+CiAgPHNjcmlwdCBzcmM9Ii8vY2RuanMuY2xvdWRmbGFyZS5jb20vYWpheC9saWJzL2hpZ2hsaWdodC5qcy8xMC4xLjIvaGlnaGxpZ2h0Lm1pbi5qcyI+PC9zY3JpcHQ+CiAgPHNjcmlwdD5obGpzLmluaXRIaWdobGlnaHRpbmdPbkxvYWQoKTs8L3NjcmlwdD4KPC9ib2R5Pgo8L2h0bWw+Cnt7IGVuZCB9fQoKe3sgdGVtcGxhdGUgImJhc2UiIC4gfX0K"
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
