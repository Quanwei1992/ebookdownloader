goqt_rcc -go main -o ebookdlform_qrc.go ebookdlform.qrc
go build -ldflags "-H windowsgui" -o ebookdl_gui.exe