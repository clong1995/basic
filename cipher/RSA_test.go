package cipher

import (
	"testing"
)

func TestRSADecrypt(t *testing.T) {
	encrypt, err := RSAEncrypt([]byte("15166077180"), []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDkebjpT6ZdlOKfH2W8j48+Qodr
+/hhnredqIF9j4CVY9b6h5E0cxy1QND8xNFhnbfByKmSP04htyKU83CZJN/kj3Ox
TDI3nH3LgXUn7+r9BQxCLoQsk3wZX2JdMNDiOJRHcabna5z4o+RL3DtSX4ojVINS
krvCb7ZLIlYJuzZAAwIDAQAB
-----END PUBLIC KEY-----`))
	if err != nil {
		t.Error(err)
		return
	}

	decrypt, err := RSADecrypt(encrypt, []byte(`-----BEGIN PRIVATE KEY-----
MIICdQIBADANBgkqhkiG9w0BAQEFAASCAl8wggJbAgEAAoGBAOR5uOlPpl2U4p8f
ZbyPjz5Ch2v7+GGet52ogX2PgJVj1vqHkTRzHLVA0PzE0WGdt8HIqZI/TiG3IpTz
cJkk3+SPc7FMMjecfcuBdSfv6v0FDEIuhCyTfBlfYl0w0OI4lEdxpudrnPij5Evc
O1JfiiNUg1KSu8JvtksiVgm7NkADAgMBAAECgYB/zec1+6wgZQxv3mxWkieauDRw
nz5NvS8RLhVhW0ieSH8VHYiIQmwop90/yAkoBcWozMquWGMoUP0zPQobYJksfsGh
9wDEoGGpTuRPo2gMtm+EphjhaSp7D6DAtfnKazSQG35filDS/08xygbnuPcnggMN
7+AtSk1pzS29vrTpYQJBAPbPNekoBGIbJ9Otbjbo2sk2WVUyW+sDgH3aI7sGgMHv
11NBMEnkeqcfIl4y8MEZ0I9vxZ4E677M17Ysce+MP9UCQQDs+7zT84WJY9QjOXLG
x/pNEXe9UhkLYVsbmtM1arzSvg6F0EGPRNCqjMlzzf2Q7WatbX4+N6eNR3+VOem3
0ER3AkBbACO8iAi1s5WHstaEYG7q6aMeiqbhjDUAMkIiX09yMmCOTebkF94xaIVf
fiDO0hnYCTov/Vh+zUBr5w9LZ8bRAkBYiOPOu1fUMDt8vWWn5eYZDMGTNSyuF70V
3w2xEyNgCCkczOTxRWA/l0FbxkVI86g8en+Ddv9dxKxhb7VlOqWZAkASGOJxupcr
1zJsTw3Q79ROi8lRGFsOCtmBtl6v3LLFqVWYC1RyiVsgbbQUKYwN7WLZGgd8nRlM
whuIDwXXFAUr
-----END PRIVATE KEY-----`))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(decrypt))
}
