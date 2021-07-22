package shell

import "strings"

var (
	// ScriptSHAs contains the default known valid SHA256 sums for bin/activate-hermit and bin/hermit.
	ScriptSHAs = []string{
		"020657425d0ba9f42fd3536d88fb3e80e85eaeae2daa1f7005b0b48dc270a084",
		"3ec9e59a260a2befeb83a94952337dddcac1fb4f7dcc1200a2964bfb336f26c3",
		"5ba24eaadfe620ad7c78a5c1f860d9845bc20077a3f9c766936485d912b75b60",
		"60c8e1787b16b6bd02c0cf6562671b0f60fb8d867b6d5140afd96bd2521e2f68",
		"6e1e6a687dc1f43c8187fb6c11b2a3ad1b1cfc93cda0b5ef307710dcfafa0dd4",
		"7a2b479e582d39826ef3e47d9930c7e8ff21275fba53efdc8204fe160742b56c",
		"04f065a430d1d99bc99f19e82a6465ab6823467d9c6b5ec3f751befa7a3b30a8",
		"57697ee9f19658d1872fc5877e2a38ba132a2df85e4416802a4c33968e00c716",
		"75abcf121df40b25cd0c7bab908c43dbf536bc6f4552a2f6e825ac90c8fff994",
		"7c64aa474afa3202305953e9b2ac96852f4bf65ddb417dee2cfa20ad58986834",
		"b42be79b29ac118ba05b8f5b6bd46faa2232db945453b1b10afc1a6e031ca068",
		"c419082d4cf1e2e9ac33382089c64b532c88d2399bae8b07c414b1d205bea74e",
		"d575eda7d5d988f6f3c233ceaa42fae61f819d863145aec7a58f4f1519db31ad",
		"ec14f88a38560d4524a8679f36fdfb2fb46ccd13bc399c3cddf3ca9f441952ec",
	}
)

func makeCommonHooks(validSHAs ...string) string {
	scriptSHAs := strings.Join(validSHAs, "|")
	return `
change_hermit_env() {
  CUR=${PWD}
  while [ "$CUR" != "/" ]; do
    if [ "${CUR}" -ef "${HERMIT_ENV}" ]; then return; fi
    if [ -f "${CUR}/bin/activate-hermit" ]; then
      if [ -n "${HERMIT_ENV+_}"  ]; then type _hermit_deactivate &>/dev/null && _hermit_deactivate; fi
      # shellcheck source=files/activate-hermit
      if [ "${CUR}" != "${DEACTIVATED_HERMIT}" ]; then
		if ! openssl sha256 "${CUR}/bin/activate-hermit" "${CUR}/bin/hermit" | awk '{print $2}' | grep -vqE '` + scriptSHAs + `'; then
			. "${CUR}/bin/activate-hermit"
		else
			echo "warning: One of ${CUR}/bin/{hermit,activate-hermit} has an unknown signature." 1>&2
			echo "         Verify that you trust this repository and run '~/bin/hermit init .'" 1>&2
		fi
      fi
      return
    fi
    CUR="$(dirname "${CUR}")"
  done
  unset DEACTIVATED_HERMIT
  if [ -n "${HERMIT_ENV+_}"  ]; then type _hermit_deactivate &>/dev/null && _hermit_deactivate; fi
}
`
}

const (
	hookStartMarker = "# Generated by Hermit; START; DO NOT EDIT."
	hookEndMarker   = "# Generated by Hermit; END; DO NOT EDIT."
)