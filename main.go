// dicma.go
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

var (
	// INTERNAL CONFIGURATION PARAMETERS FOR PASSWORD MODE:
	amount_of_sufixs_used_light_mode  = 200
	amount_of_prefixs_used_light_mode = 50

	amount_of_sufixs_used  = len(BASIC_SUFIXS)
	amount_of_prefixs_used = 200

	amount_of_numericpat_used = 516
	amount_of_symbolpat_used  = 100

	// General variables
	LIGHT_MODE          = false
	FULL_MODE           = false
	VERBOSE             = true
	OUTPUT_FILE_BOOLEAN = false
	MASSIVE_MODE        = false

	NEIGHBORS_AMMOUNT = 20

	// Patterns (copiados del python)
	BASIC_SUFIXS = []string{
		"1", "2", "123", "12", "3", "7", "13", "5", "4", "11", "!", "07", "23", "22", "01", "21", "8", "14", "10", "08", "6", "06", "9", "15", "16", "69", "18", "17", "24", "05", ".", "09", "88", "19", "25", "20", "03", "0", "04", "27", "89", "02", "99", "26", "101", "77", "1234", "28", "33", "00", "2007", "92", "87", "93", "2006", "*", "29", "94", "90", "2008", "91", "95", "86", "55", "30", "666", "143", "31", "96", "85", "44", "32", "007", "34", "4ever", "84", "45", "#1", "2005", "78", "98", "66", "83", "82", "97", "1994", "79", "100", "1992", "1993", "81", "4life", "4eva", "12345", "1995", "1991", "777", "1990", "420", "76", "56", "111", "1989", "2000", "1987", "321", "2009", "67", "80", "75", "2004", "1996", "42", "1988", "35", "74", "72", "1986", "1985", "001", "36", "73", "54", "2003", "456", "50", "!!", "333", "68", "@", "1984", "64", "37", "65", "40", "71", "2002", "4u", "123456", "43", "555", "911", "1997", "52", "?", "999", "1983", "1982", "47", "$", "57", "2001", "41", "@hotmail.com", "1980", "38", "4me", "1981", "63", "46", "70", "2010", "58", "48", "222", "51", "62", "1979", "121", "619", "53", "789", "59", "39", "112", "1998", "1978", "000", "888", "49", "**", "247", "234", "61", "60", "..", "213", "1977", "1212", "200", "159", "1999", "@yahoo.com", "1976", "182", "786", "1!", ".com", "<3", "1975", ")", "125", "187", "214", "2k7", "147", "...", "1974", ".1", "102", "1973", "!!!", "4e", "123456789", "411", "212", "6969", "1010", "305", "124", "012", "1972", "345", "360", "1969", "009", "316", "1970", "313", "4lyf", "711", "210", "808", "122", "1313", "987", "311", "120", "#", "1971", "444", "369", "1000", "008", "500", "2012", "323", "2011", "211", "1111", "246", "215", ",", "300", "135", "567", "117", "003", "103", "113", "1968", "1122", "225", "713", "132", "209", "510", "002", "520", "1221", "310", "223", "105", "~", "127", "100pre", "110", "1967", "818", "011", "1213", "1012", "202", "312", "128", "109", "+", "108", "126", "098", "1966", "2k6", "1965", "909", "1.", "118", "1020", "3000", "115", "678", "714", "415", "1964", "013", "1123", "107", "718", "129", "]", "831", "010", "2121", "4l", "005", "145", "314", "224", "104", "2525", "131", ";", "221", "626", "1314", "2020", "357", "4lyfe", "114", "318", "1210", "1223", "916", "1023", "1230", "106", "900", "216", "006", "116", "1011", "2323", "201", "2468", "4321", "119", "504", "3r", "1963", "1214", "217", "315", "2u", "512", "1224", "412", "0123", "1231", "1022", "1121", "1013", "1001", "258", "890", "134", "150", "220", "421", "1021", "1220", "1024", "1218", "1215", "317", "199", "521", "1216", "***", "303", "515", "004", "707", "$$", "206", "1014", "219", "1025", "1211", "1962", "218", "320", "1960", "1217", "413", "4444", "727", "1205", "1225", "`", "912", "812", "410", "0n", "813", "227", "2k", "014", "100%", "1206", "231", "1228", "505", "1029", "7777", "0000", "1017", "1015", "1107", "133", "1112", "325", "130", "513", "4LIFE", "324", "@1", "1016", "511", "1018", "423", "180", "1026", "021", "1125", "1031", "1028", "'", "1203", "1204", "1104", "023", "2k8", "612", "400", "721", "805", "1105", "2112", "1124", "1120", "4you", "525", "1207", "228", "205", "913", "2222", "2013", "#2", "1103", "1202", "1961", "1027", "617", "1106", "1019", "1127", "717", "1227", "156", "250", "4evr", "1208", "712", "600", "1219", "254", "1414", "017", "226", "408", "925", "723", "322", "1959", "137", "817", "319", "2424", "151", "809", "910", "232", "/", "203", "1229", "414", "326", "327", "616", "1515", "1226", "1030", "516", "1*", "1201", "2x", "365", "144", "1129", "123!", "623", "1209", "523", "611", "1126", "1005", "9999", "1128", "=", "1004", "!1", "1102", "710", "1002", "654", "1101", "613", "419", "1108", "422", "921", "198", "1130", "425", "123.", "923", "915", "169", "615", "015", "416", "#3", "951", "610", "141", "190", "816", "328", "1958", "811", "720", "016", "828", "426", "702", "919", "191", "1a", "522", "1007", "501", "018", "823", "189", "715", "8888", "918", "54321", "716", "256", "1003", "424", "1234567", "963", "138", "@123", "207", "427", "2345", "787", "5150", "852", "242", "914", "621", "331", "4EVER", "741", "637", "614", ".123", "562", "252", "920", "301", "1109", "1222", "153", "330", "1e", "417", "1006", "530", "800", "0506", "1957", "*1", "245", "1956", "517", "622", "155", "235", "753", "514", "329", ":)", "1113", "0101", "157", "821", "722", "810", "5678", "429", "2526", "724", "527", "024", "069", "524", "418", "0607", "208", "922", "624", "142", "168", "152", "2me", "0808", "1008", "929", "561", "2b", "700", "1717", "928", "815", "233", "518", "136", "177", "0708", "625", "0102", "725", "901", "022", "822", "917", "0909", "519", "334", "526", "248", "183", "731", "927", "531", "1117", "167", "602", "1114", "1818", "814", "5555", "529", "019", "620", "1n", "1616", "1955", "824", "618", "%", "1009", "192", "719", "503", "924", "1919", "229", "1st", "@aol.com", "926", "'s", "820", "430", "0707", "528", "188", "181", "4EVA", "989", "432", "1412", "204", "404", "1110", "0406", "304", "428", "5000", "0405", "825", "0809", "3s", "??", "628", "0204", "281", "2212", "3d", "0205", "197", "0202", "1115", "1954", "726", "728", "237", "0505", "757", "1116", "302", "140", "954", "559", "601", "747", "1516", "405", "730", "0305", "243", "146", "154", "509", "729", "627", "0306", "1415", "\"", "350", "629", "1324", "696", "0507", "0711", "0214", "253", "269", "407", "0407", "0304", "306", "4U", "650", "2310", "171", "236", "2529", "630", "1905", "239", "2528", "450", "1s", "2527", "139", "025", "195", "&", "0303", "123*", "1118", "454", "0606", "2530", "0203", "876", "930", "178", "\\'", "2210", "0107", "0404", "027", "0408", "1907", "819", "0812", "409", "158", "908", "543", "401", "230", "0210", "1312", "0308", "307", "1love", "240", "0412", "1369", "1310", "1953", "826", "829", "2!", "165", "2105", "161", "606", "1234567890", "160", "148", "904", "0311", "830", "0307", "0103", "2211", "2014", "193", ".12", "0110", "308", "827", "1307", "0420", "0206", "1119", "0207", "@@", "0911", "0212", "149", "990", "0105", "2123", ".2", "255", "309", "(L)", "241", "3333", "801", "0912", "185", "1432", "0104", "6666", "0987", "175", "163", "502", "1950", "671", "0106", "0208", "765", "251", "337", "1908", "170", "445", "166", "028", "173", "257", "3y", "031", "2107", "1888", "2524", "123123", "196", "&me", "1402", "0608", "2311", "1331", "!@#", "1903", "1305", "0312", "1235", "1@", "956", "174", "2312", "0509", "275", "2523", "0508", "172", "0512", "1311", "1306", "1952", "1200", "90210", "1512", "1408", "1233", "026", "1411", "2727", "1920", "978", "184", "5683", "1690", "0810", "2510", "244", "0709", "565", "2410", "667", "238", "090", "0108", "164", "12345678", "55555", "4u2", "0411", "1912", "#7", "3n", "2412", "2205", "1812", "1904", "0811", "2103", "402", "0612", "0309", "973", "0907", "336", "14344", "1410", "123321", "1906", "1323", "289", "903", "\\", "0211", "0209", "1405", "2106", "2828", "0712", "803", "2580", "2626", "267", "0321", "1308", "1718", ">", "2207", "1304", "262", "4ev", "077", "179", "343", "2531", "2110", "3030", "1508", "1407", "12!", "2203", "162", "0x", "585", "0906", "176", "1821", "1606", "2303",
	}

	BASIC_PREFIXS = []string{
		"1", "2", "4", "123", "3", "*", "12", "7", "5", "8", "6", "13", "11", "143", "9", "19", "@", "(", "10", "14", "0", "22", "23", "21", "$", "!", "#1", "20", "15", "24", "18", "17", "16", "1234", "69", ".", "25", "01", "27", "i", "07", "26", "**", "28", "06", "29", "08", "30", "~", "00", "05", "123456", "88", "99", "12345", "02", "03", "100", "04", "31", "09", "[", "666", "33", "77", "#", "sk8", "50", "ms.", "mr.", "123456789", "1994", "i<3", "100%", "mz.", "<", "89", "l0", "44", "1992", "1993", "420", "32", "1995", "55", "2007", "101", "34", "321", "1991", "007", "92", ",", "m1", "2006", "87", "95", "1989", "\"", "<3", "1996", "1990", "98", "no1", "93", "1987", "777", "m0", "94", "c0", "90", "66", "2008", "45", ";", "r0", "91", "96", "il0", "86", "1988", "my1", "ih8", "78", "97", "my", "h0", "..", "+", "!!", "111", "p0", "l1", "j0", "56", "b1", "m3", "85", "1986", "1985", "52", "d3", "35", "42", "k1", "/", "s3", "2005", "84", "40", "79", "`", "82", "619", "***", "s0", "54", "p1", "b3", "76", "te", "36", "74", "67", "a1", "g0", "hi5", "=", "83", "d0", "1997", "911", "333", "1984", "b0", "49", "s1", "456", "d1", "68", "555", "m@", "80", "72", "43", "j.", "?", "$$", "...", "j3", "789", "h3", "999", "187", "37", "1983", "75", "57", "n0", "41",
	}

	NUMERIC_PATTERNS = []string{
		"1", "2", "4", "3", "123", "7", "12", "5", "0", "8", "13", "6", "9", "11", "23", "22", "10", "14", "07", "21", "01", "15", "08", "06", "16", "18", "69", "17", "24", "05", "19", "09", "20", "25", "88", "03", "00", "27", "04", "02", "33", "89", "26", "99", "1234", "28", "77", "101", "92", "2007", "93", "87", "29", "94", "2006", "143", "90", "91", "30", "95", "2008", "55", "31", "100", "86", "666", "34", "44", "32", "96", "85", "45", "007", "84", "98", "1994", "78", "66", "2005", "1992", "1993", "83", "82", "12345", "97", "1991", "1995", "79", "1990", "81", "1989", "56", "777", "1987", "76", "420", "321", "111", "35", "67", "1996", "80", "2000", "1988", "42", "75", "123456", "50", "2009", "74", "2004", "1986", "72", "54", "36", "456", "1985", "73", "43", "1984", "64", "2003", "68", "37", "40", "333", "001", "52", "65", "47", "41", "71", "1997", "57", "1983", "555", "2002", "999", "1982", "911", "38", "1980", "2001", "46", "63", "70", "1981", "48", "51", "53", "58", "789", "62", "39", "2010", "59", "619", "49", "1979", "222", "121", "000", "1998", "112", "1978", "123456789", "60", "61", "888", "234", "159", "1212", "247", "200", "1977", "213", "786", "1999", "125", "182", "1976", "187", "1975", "147", "214", "1974", "1973", "102", "305", "1010", "6969", "124", "1972", "212", "360", "987", "1000", "411", "345", "1969", "012", "808", "313", "311", "210", "1313", "1970", "369", "500", "120", "300", "316", "711", "122", "1971", "009", "1111", "444", "323", "520", "135", "2012", "246", "211", "2011", "1122", "567", "215", "103", "008", "113", "1968", "117", "225", "510", "310", "209", "1221", "110", "132", "2468", "713", "105", "003", "127", "1213", "1967", "312", "002", "223", "1012", "109", "128", "011", "818", "1966", "108", "126", "098", "202", "118", "2525", "1020", "1123", "104", "831", "1965", "115", "357", "010", "2121", "107", "714", "1314", "145", "909", "415", "900", "129", "678", "314", "224", "1964", "318", "2020", "1210", "114", "131", "2323", "718", "3000", "504", "013", "106", "4321", "1223", "1230", "221", "626", "1023", "201", "116", "005", "1011", "150", "916", "216", "512", "119", "0123", "1214", "515", "006", "412", "1963", "1224", "315", "258", "0000", "134", "217", "220", "812", "1121", "199", "890", "1001", "1231", "1024", "521", "1021", "400", "1022", "1013", "4444", "303", "1215", "1220", "421", "250", "1234567", "1216", "206", "320", "1218", "317", "707", "325", "218", "1025", "1962", "7777", "413", "1211", "1014", "219", "410", "133", "1225", "130", "600", "1029", "1205", "1960", "813", "1217", "004", "1206", "1015", "231", "1228", "1112", "912", "227", "180", "1107", "505", "654", "727", "014", "1017", "324", "1120", "1204", "1031", "513", "2222", "1016", "423", "1028", "2112", "1026", "805", "205", "1203", "156", "1018", "963", "1105", "1125", "1414", "228", "511", "1124", "021", "408", "612", "1104", "525", "2424", "1208", "232", "254", "151", "1202", "1027", "809", "319", "1103", "144", "1207", "1515", "1106", "1127", "322", "721", "913", "203", "414", "226", "54321", "137", "023", "2013", "1005", "1019", "800", "1219", "1227", "910", "1961", "326", "717", "925", "617", "951", "1229", "723", "1030", "611", "1129", "1209", "700", "817", "1959", "1226", "327", "712", "365", "1004", "710", "616", "017", "9999", "1201", "1102", "5150", "1128", "852", "198", "425", "191", "1002", "523", "623", "426", "141", "350", "1101", "516", "256", "753", "1126", "422", "419", "8888", "741", "155", "1108", "1130", "169", "787", "190", "1234567890", "501", "252", "242", "613", "915", "610", "2014", "2015", "2016", "2017", "2018", "2019", "2021", "2021", "2023", "2024", "2025", "2026", "2027", "2028", "2029", "2030",
	}

	SYMBOLIC_PATTERNS = []string{
		".", "-", "!", "@", "*", "/", "#", "&", ",", "$", "+", "=", "?", "(", ")", "**", "!!", ";", "<", "..", "'", "]", "%", "\"", "~", "...", "[", "`", "=\"", "\\'", ":", "!!!", "$$", "***", "^", "--", "@@", "//", "///", ">", "++", "??", "!@", "\\", ":)", "ื", "://", "!@#", "##", "\\\\\\'", ".,", "{", "\\\\", ",.", "}", "$$$", "><", ",,", "()", "ั", "้", "/*", "^^", "ุ", "@#", "ึ", "+-", "&&", "???", "@@@", "*/", "ิ", "|", ";;", "ี", "****", "==", "@!", "....", "!*", "[]", ",./", "---", "=]", "/*-", "@$", "=)", "!!!!", ".-", "#!", "~~", ",]", "´", "!@#$", "-=", "*-", ")(", "+++", "))", "?!", "=-",
	}
)

func printBanner() {
	ascii := `
                                 .-:                                  
                               -*%%%#+.                               
                             =#%%%%%%%%+.                             
                           .*%%%%%%%%%%%#-                            
                          :#%%%%%%%%%%%%%#=                           
                         :#*%%%#=:..=*%%%*#=                          
                         ##%%+:       .=#@#%:                         
                        =%%*.           .=%%#                         
                        #%=               :#%:                        
                       .@= .              .:%=                        
                        #.:-             :- *:                        
                       ...:%=           :%+ :.   -                    
                      :%+. +%+         -%#: -#: *:                    
                    ..::--:..=+:     .++:.:+=.:#= =                   
                   .-=+++++=-::::   .:.  .. .+#- =-                   
                  -+=-:...::-=+++=:.     .-*#+  *+ :                  
                    :-=++++=-:.  ... .:=#%*-  -#- -.                  
                 .=+=-:..........:=+###+-  .=#+. =:.                  
                 :.  .:::.  .-+#%#*=:   .-*#+. -*:.:                  
                   :-:.  .-#%#+-.   .-+#*+:  -*= :-                   
                 .:.   .+%#+:    :=+*+-.  :=+=. --                    
                      =%#-    :=++-:   .-+=:  :-.                     
                     +%-   .-==-.   .-==:   :-:                       
                    =*.   :=-.   .:--:   .:-:                         
                   :+    --.   .:-:    .::.                           
                   :    ::    .:.    .:.                              
                       :.    :.    ...                                
                      ..    :.    ..                                  
                                                                      `
	fmt.Println(ascii)
	fmt.Println("Welcome to DICMA. The Dictionary Maker: \n")
}

func verbosePrint(s string) {
	if VERBOSE {
		fmt.Println(s)
	}
}

func detectIfFileOrNot(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func systemDetection() string {
	if runtime.GOOS == "windows" {
		return "windows"
	}
	return "linux"
}

func getTotalRAM() (float64, error) {
	// Attempt platform-specific methods. Keep it simple and approximate.
	switch runtime.GOOS {
	case "linux":
		// read /proc/meminfo
		content, err := ioutil.ReadFile("/proc/meminfo")
		if err != nil {
			return 0, err
		}
		re := regexp.MustCompile(`MemTotal:\s+(\d+)`)
		m := re.FindSubmatch(content)
		if len(m) >= 2 {
			// kB to GB
			var kb int64
			fmt.Sscanf(string(m[1]), "%d", &kb)
			gb := float64(kb) / (1024.0 * 1024.0)
			return round(gb, 2), nil
		}
	case "darwin":
		// sysctl hw.memsize
		out, err := exec.Command("sysctl", "hw.memsize").Output()
		if err != nil {
			return 0, err
		}
		parts := strings.Split(string(out), ":")
		if len(parts) >= 2 {
			var bytesVal float64
			fmt.Sscanf(strings.TrimSpace(parts[1]), "%f", &bytesVal)
			gb := bytesVal / (1024.0 * 1024.0 * 1024.0)
			return round(gb, 2), nil
		}
	case "windows":
		// use wmic if available
		out, err := exec.Command("wmic", "computersystem", "get", "totalphysicalmemory").Output()
		if err == nil {
			lines := strings.Fields(string(out))
			if len(lines) >= 2 {
				var bytesVal float64
				fmt.Sscanf(lines[1], "%f", &bytesVal)
				gb := bytesVal / (1024.0 * 1024.0 * 1024.0)
				return round(gb, 2), nil
			}
		}
	}
	return 0, fmt.Errorf("unsupported or cannot detect")
}

func round(x float64, prec int) float64 {
	p := 1.0
	for i := 0; i < prec; i++ {
		p *= 10
	}
	return float64(int64(x*p+0.5)) / p
}

func isAValidFile(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	buf := make([]byte, 8192)
	_, err = f.Read(buf)
	if err != nil && err != io.EOF {
		return false
	}
	// Python's check attempted to decode as latin-1, which always succeeds for bytes,
	// so we'll consider the file valid if readable.
	return true
}

func saveListToFile(list []string, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, it := range list {
		_, _ = w.WriteString(it + "\n")
	}
	w.Flush()
	verbosePrint("[+] Dictionary saved successfully to: " + filename)
	return nil
}

func removeAccents(s string) string {
	// Normalize NFD and remove marks
	t := norm.NFD.String(s)
	var b strings.Builder
	for _, r := range t {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func askForYesOrNo(question string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s (yes/no): ", question)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(strings.ToLower(text))
		if text == "yes" || text == "y" {
			return true
		}
		if text == "no" || text == "n" {
			return false
		}
		fmt.Println("Please, answer 'yes' or 'no'.")
	}
}

func generateUsernames(personName string) []string {
	parts := strings.Fields(strings.TrimSpace(personName))
	if len(parts) <= 1 {
		return parts
	}
	if len(parts) >= 3 {
		fmt.Fprintf(os.Stderr, "[!] Unsuported names of 3 words -> %s\n", personName)
		os.Exit(1)
	}
	name := parts[0]
	surname := parts[1]
	firstLetName := ""
	firstLetSurname := ""
	if len(name) > 0 {
		firstLetName = string([]rune(name)[0])
	}
	if len(surname) > 0 {
		firstLetSurname = string([]rune(surname)[0])
	}

	combinations := []string{
		fmt.Sprintf("%s %s", name, surname),
		fmt.Sprintf("%s%s", name, surname),
		fmt.Sprintf("%s%s", firstLetName, surname),
		fmt.Sprintf("%s%s", name, firstLetSurname),
		fmt.Sprintf("%s.%s", name, surname),
		fmt.Sprintf("%s.%s", firstLetName, surname),
		fmt.Sprintf("%s.%s", name, firstLetSurname),
		fmt.Sprintf("%s_%s", name, surname),
		fmt.Sprintf("%s_%s", firstLetName, surname),
		fmt.Sprintf("%s_%s", name, firstLetSurname),
		fmt.Sprintf("%s-%s", name, surname),
		fmt.Sprintf("%s-%s", firstLetName, surname),
		fmt.Sprintf("%s-%s", name, firstLetSurname),
	}
	return combinations
}

func normalizeList(input_ string) []string {
	if detectIfFileOrNot(input_) {
		// read file lines
		file, err := os.Open(input_)
		if err != nil {
			return nil
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		words := []string{}
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				words = append(words, line)
			}
		}
		return words
	} else {
		parts := strings.Split(input_, ",")
		res := []string{}
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				res = append(res, p)
			}
		}
		return res
	}
}

func processFileUser(fileName, outputFileName string) {
	if !OUTPUT_FILE_BOOLEAN {
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
			return
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			combinations := generateUsernames(line)
			for _, e := range combinations {
				fmt.Println(e)
			}
		}
	} else {
		completList := []string{}
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
			return
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			combinations := generateUsernames(line)
			completList = append(completList, combinations...)
		}
		err = saveListToFile(completList, outputFileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving file: %v\n", err)
		}
	}
}

func processInputUser(input_, outputFileName string) {
	parts := strings.Split(input_, ",")
	peopleList := []string{}
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			peopleList = append(peopleList, p)
		}
	}
	finalList := []string{}
	if !OUTPUT_FILE_BOOLEAN {
		for _, n := range peopleList {
			finalList = append(finalList, generateUsernames(n)...)
		}
		for _, e := range finalList {
			fmt.Println(e)
		}
	} else {
		for _, n := range peopleList {
			finalList = append(finalList, generateUsernames(n)...)
		}
		err := saveListToFile(finalList, outputFileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving file: %v\n", err)
		}
	}
}

func processPasswd(wordsList []string, outputFileName string) {
	verbosePrint("[+] Creating dictionary.")
	// massive mode checks
	if len(wordsList) >= 2 && FULL_MODE {
		massiveMode(wordsList, outputFileName)
		return
	}
	if len(wordsList) >= 10 && !LIGHT_MODE {
		massiveMode(wordsList, outputFileName)
		return
	}
	if len(wordsList) >= 100 {
		massiveMode(wordsList, outputFileName)
		return
	}

	finalList := []string{}
	if !OUTPUT_FILE_BOOLEAN {
		for _, w := range wordsList {
			finalList = append(finalList, generatePasswordList(w, FULL_MODE, LIGHT_MODE)...)
		}
		for _, e := range finalList {
			fmt.Println(e)
		}
	} else {
		for _, w := range wordsList {
			finalList = append(finalList, generatePasswordList(w, FULL_MODE, LIGHT_MODE)...)
		}
		err := saveListToFile(finalList, outputFileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving file: %v\n", err)
		}
	}
}

// findNeighbours - linux path uses python helper; windows path uses fasttext.exe batch if present.
func findNeighboursBatchWindows(modelPath string, words []string, number int) []string {
	// replicate the behavior: call fasttext.exe nn modelPath number, feed words via stdin
	cmd := exec.Command("fasttext.exe", "nn", modelPath, fmt.Sprintf("%d", number))
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] Error creating stdin pipe: %v\n", err)
		return []string{}
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "[!] Error starting fasttext.exe: %v\n", err)
		return []string{}
	}
	// write input words
	for _, w := range words {
		io.WriteString(stdin, w+"\n")
	}
	stdin.Close()
	err = cmd.Wait()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] Error from fasttext.exe: %v -> %s\n", err, stderr.String())
		return []string{}
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	// parse output similarly to python helper
	outputMap := map[string][]string{}
	wordIndex := 0
	currentNeighbors := []string{}
	for _, line := range lines {
		if strings.HasPrefix(line, "Query word?") {
			if len(currentNeighbors) > 0 && wordIndex < len(words) {
				w := words[wordIndex]
				outputMap[w] = cleanNeighbors(currentNeighbors)
				currentNeighbors = []string{}
				wordIndex++
			}
		} else {
			currentNeighbors = append(currentNeighbors, line)
		}
	}
	if len(currentNeighbors) > 0 && wordIndex < len(words) {
		w := words[wordIndex]
		outputMap[w] = cleanNeighbors(currentNeighbors)
	}
	unified := map[string]struct{}{}
	for w, neigh := range outputMap {
		unified[w] = struct{}{}
		for _, n := range neigh {
			unified[n] = struct{}{}
		}
	}
	final := []string{}
	for k := range unified {
		final = append(final, k)
	}
	return final
}

func cleanNeighbors(lines []string) []string {
	neigh := []string{}
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 1 {
			neighbor := parts[0]
			neighbor = strings.TrimRight(neighbor, "-")
			neighbor = strings.ToLower(neighbor)
			if strings.Contains(neighbor, ".") {
				continue
			}
			neigh = append(neigh, neighbor)
		}
	}
	// unique
	set := map[string]struct{}{}
	out := []string{}
	for _, w := range neigh {
		if _, ok := set[w]; !ok {
			set[w] = struct{}{}
			out = append(out, w)
		}
	}
	return out
}

func mlProcessPwd(list_ []string, mlModel string, numberNeighbours int) []string {
	systemOS := systemDetection()
	if systemOS == "linux" {
		// We'll try to run a small python helper that loads fasttext and prints neighbours as JSON.
		// This requires python and fasttext installed (same requirement as original).
		py := `
		import sys, json
		try:
			import fasttext
		except Exception as e:
			print("ERROR:fasttext_not_loaded", file=sys.stderr)
			sys.exit(1)
		model = fasttext.load_model(sys.argv[1])
		k = int(sys.argv[2])
		words = []
		for line in sys.stdin:
			w = line.strip()
			if not w: continue
			words.append(w)
		res = {}
		for w in words:
			n = model.get_nearest_neighbors(w, k)
			neighs = [t for _, t in n if '.' not in t]
			# filter uppercase in position >0
			def_valid = []
			for t in neighs:
				ok = True
				for ch in t[1:]:
					if ch.isupper():
						ok = False
						break
				if ok:
					def_valid.append(t)
			res[w] = def_valid
		print(json.dumps(res))
		`
		cmd := exec.Command("python3", "-c", py, mlModel, fmt.Sprintf("%d", numberNeighbours))
		in := bytes.Buffer{}
		for _, w := range list_ {
			in.WriteString(w + "\n")
		}
		cmd.Stdin = &in
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "[-] Unable to load fasttext module or python error: %v - %s\n", err, stderr.String())
			os.Exit(1)
		}
		var parsed map[string][]string
		if err := json.Unmarshal(out.Bytes(), &parsed); err != nil {
			fmt.Fprintf(os.Stderr, "[-] Error parsing fasttext output: %v\n", err)
			os.Exit(1)
		}
		wordsList := []string{}
		for _, w := range list_ {
			verbosePrint("[+] Processing -> " + w)
			if neigh, ok := parsed[w]; ok {
				wordsList = append(wordsList, neigh...)
			}
		}
		verbosePrint(fmt.Sprintf("[+] Neighbors found successfully, %d Words in total.", len(wordsList)))
		return wordsList
	}
	if systemOS == "windows" {
		verbosePrint("[!] Warning, WINDOWS OS detected. \"fasttext\" library can not be compiled in windows. We are going to use fasttext.exe precompiled binary.")
		if !detectIfFileOrNot("fasttext.exe") {
			verbosePrint("[!] fasttext.exe not found in present directory...")
			verbosePrint("[!] We need to download an external fasttext.exe from -> https://github.com/sigmeta/fastText-Windows/releases/")
			if askForYesOrNo("Do you accept ? (answer yes or no)") {
				// try to download using powershell or curl? We'll try to use powershell to download.
				if !isWritable(".") {
					fmt.Fprintln(os.Stderr, "[-] We don't have write permision in this folder. Move to a folder where you can write, or download fasttext.exe by yourself and move it to this folder.")
					os.Exit(1)
				}
				url := "https://github.com/sigmeta/fastText-Windows/releases/download/0.9.2/fasttext.exe"
				verbosePrint("[!] Trying to download " + url)
				// attempt with powershell Invoke-WebRequest or curl
				var err error
				if detectIfFileOrNot("curl.exe") {
					err = exec.Command("curl.exe", "-L", "-o", "fasttext.exe", url).Run()
				} else {
					// powershell
					ps := fmt.Sprintf(`(New-Object System.Net.WebClient).DownloadFile("%s","fasttext.exe")`, url)
					err = exec.Command("powershell", "-Command", ps).Run()
				}
				if err != nil {
					fmt.Fprintf(os.Stderr, "[-] Error downloading fasttext.exe: %v\n", err)
					os.Exit(1)
				}
				fmt.Println("[+] Binary downloaded successfully")
			} else {
				fmt.Fprintln(os.Stderr, "[-] Sorry but without this binary we won't be able to use fasttext model in Windows. Try to use dicma in linux, macos, or get this binary.")
				os.Exit(1)
			}
		}
		verbosePrint("[+] Looking up to " + fmt.Sprintf("%d", numberNeighbours) + " nearest neighbours for each word.")
		list := findNeighboursBatchWindows(mlModel, list_, numberNeighbours)
		verbosePrint("[+] Neighbors found successfully.")
		return list
	}
	return []string{}
}

func isWritable(path string) bool {
	test := filepath.Join(path, ".dicma_write_test")
	err := ioutil.WriteFile(test, []byte("x"), 0644)
	if err != nil {
		return false
	}
	os.Remove(test)
	return true
}

func workerGenerate(word string, FULL_MODE, LIGHT_MODE bool, outChan chan<- string, batchSize int, wg *sync.WaitGroup) {
    defer wg.Done()
    batch := []string{}
    for _, line := range generatePasswordList(word, FULL_MODE, LIGHT_MODE) {
        batch = append(batch, line)
        if len(batch) >= batchSize {
            outChan <- strings.Join(batch, "\n")
            batch = []string{}
        }
    }
    if len(batch) > 0 {
        outChan <- strings.Join(batch, "\n")
    }
}

func massiveMode(list_ []string, outputFileName string) {
    verbosePrint("[!] Massive mode ENABLED")
    VERBOSE = true
    if !OUTPUT_FILE_BOOLEAN {
        outputFileName = "output.txt"
        verbosePrint("[!] Output file required for the massive mode. Saving results to -> " + outputFileName)
    }

    cpuCores := runtime.NumCPU()
    if FULL_MODE {
        ram, err := getTotalRAM()
        if err == nil {
            verbosePrint(fmt.Sprintf("[i] %.2f GB RAM detected.", ram))
            if ram <= 31 {
                cpuCores = cpuCores / 2
                if cpuCores == 0 {
                    cpuCores = 1
                }
            }
            if ram <= 15 {
                verbosePrint("[!] WARNING! Full mode + multicore could saturate RAM, crash, and go slower than expected. \n[!] We will use just few cores... Go for a coffe.")
                cpuCores = cpuCores / 4
                if cpuCores == 0 {
                    cpuCores = 1
                }
            }
        }
    }
    verbosePrint(fmt.Sprintf("[+] Using %d CPU cores", cpuCores))

    outFile, err := os.Create(outputFileName)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
        return
    }
    defer outFile.Close()
    writer := bufio.NewWriter(outFile)

    outChan := make(chan string, 500)
    var wg sync.WaitGroup

    // 1) Lanzar writer goroutine primero
    doneWriter := make(chan struct{})
    go func() {
        for chunk := range outChan {
            writer.WriteString(chunk + "\n")
        }
        writer.Flush()
        close(doneWriter)
    }()

    // 2) Lanzar workers
    sem := make(chan struct{}, cpuCores)
    for _, w := range list_ {
        wg.Add(1)
        sem <- struct{}{}
        go func(word string) {
            defer func() { <-sem }()
            workerGenerate(word, FULL_MODE, LIGHT_MODE, outChan, 200, &wg)
        }(w)
    }

    // 3) Esperar workers
    wg.Wait()
    close(outChan)   // cerrar canal cuando ya no habrá más datos

    // 4) Esperar al writer
    <-doneWriter

    fmt.Printf("\n[+] Finished. Output saved to %s\n", outputFileName)
}

func generatePasswordList(word string, FULL_MODE, LIGHT_MODE bool) []string {
    localAmountSuf := amount_of_sufixs_used
    localAmountPref := amount_of_prefixs_used
    if LIGHT_MODE {
        localAmountSuf = amount_of_sufixs_used_light_mode
        localAmountPref = amount_of_prefixs_used_light_mode
    }

    // Pre-calcular tamaño estimado para optimizar memoria
    estimatedSize := 1000
    if FULL_MODE {
        estimatedSize = 50000
    } else if LIGHT_MODE {
        estimatedSize = 200
    }
    
    nonRepeatedList := make([]string, 0, estimatedSize)
    basicPattern := make([]string, 0, 20)
    
    w := strings.ToLower(strings.TrimSpace(word))
    wNoPunct := removeAccents(w)
    
    // Patrones básicos (EXACTAMENTE como en la versión original)
    basicPattern = append(basicPattern,
        w,
        strings.ToUpper(w),
        strings.Title(w),
        wNoPunct,
        strings.ToUpper(wNoPunct),
        strings.Title(wNoPunct),
    )
    
    // 🚫 RESTAURAR LA LÓGICA ORIGINAL COMPLETA sin cambios
    
    // Transformación a → @ (ORIGINAL)
    transform1 := []string{}
    for _, item := range basicPattern {
        item2 := strings.ReplaceAll(item, "a", "@")
        item2 = strings.ReplaceAll(item2, "A", "@")
        transform1 = append(transform1, item2)
    }
    basicPattern = append(basicPattern, transform1...)

    // Transformación o → 0 (ORIGINAL)
    transform2 := []string{}
    for _, item := range basicPattern {
        item2 := strings.ReplaceAll(item, "o", "0")
        item2 = strings.ReplaceAll(item2, "O", "0")
        transform2 = append(transform2, item2)
    }
    basicPattern = append(basicPattern, transform2...)

    if !LIGHT_MODE {
        // Transformación e → € (ORIGINAL)
        transform3 := []string{}
        for _, item := range basicPattern {
            item2 := strings.ReplaceAll(item, "e", "€")
            item2 = strings.ReplaceAll(item2, "E", "€")
            transform3 = append(transform3, item2)
        }
        basicPattern = append(basicPattern, transform3...)

        // Transformación e → 3 (ORIGINAL)
        transform4 := []string{}
        for _, item := range basicPattern {
            item2 := strings.ReplaceAll(item, "e", "3")
            item2 = strings.ReplaceAll(item2, "E", "3")
            transform4 = append(transform4, item2)
        }
        basicPattern = append(basicPattern, transform4...)

        // Transformación s → $ (ORIGINAL)
        transform5 := []string{}
        for _, item := range basicPattern {
            item2 := strings.ReplaceAll(item, "s", "$")
            item2 = strings.ReplaceAll(item2, "S", "$")
            transform5 = append(transform5, item2)
        }
        basicPattern = append(basicPattern, transform5...)

        // Transformación l → 1 (ORIGINAL)
        transform6 := []string{}
        for _, item := range basicPattern {
            item2 := strings.ReplaceAll(item, "l", "1")
            item2 = strings.ReplaceAll(item2, "L", "1")
            transform6 = append(transform6, item2)
        }
        basicPattern = append(basicPattern, transform6...)
    }

    // Hacer unique de basicPattern aquí (como en el original)
    basicPattern = uniquePreserve(basicPattern)
    nonRepeatedList = append(nonRepeatedList, basicPattern...)

    // word + suffixs (ORIGINAL)
    limitSuf := localAmountSuf
    if limitSuf > len(BASIC_SUFIXS) {
        limitSuf = len(BASIC_SUFIXS)
    }
    for _, a := range basicPattern {
        for _, b := range BASIC_SUFIXS[:limitSuf] {
            nonRepeatedList = append(nonRepeatedList, a+b)
        }
    }

    // prefix + word (ORIGINAL)
    limitPref := localAmountPref
    if limitPref > len(BASIC_PREFIXS) {
        limitPref = len(BASIC_PREFIXS)
    }
    for _, a := range BASIC_PREFIXS[:limitPref] {
        for _, b := range basicPattern {
            nonRepeatedList = append(nonRepeatedList, a+b)
        }
    }

    // prefix + word + suffix (ORIGINAL)
    prefLimit := limitPref / 3
    sufLimit := limitSuf / 3
    if prefLimit < 0 {
        prefLimit = 0
    }
    if sufLimit < 0 {
        sufLimit = 0
    }
    for _, a := range BASIC_PREFIXS[:prefLimit] {
        for _, b := range basicPattern {
            for _, c := range BASIC_SUFIXS[:sufLimit] {
                nonRepeatedList = append(nonRepeatedList, a+b+c)
            }
        }
    }

    // EXTENDED MODE (ORIGINAL)
    if FULL_MODE {
        numLimit := amount_of_numericpat_used
        if numLimit > len(NUMERIC_PATTERNS) {
            numLimit = len(NUMERIC_PATTERNS)
        }
        symLimit := amount_of_symbolpat_used
        if symLimit > len(SYMBOLIC_PATTERNS) {
            symLimit = len(SYMBOLIC_PATTERNS)
        }

        // word + number
        for _, a := range basicPattern {
            for _, b := range NUMERIC_PATTERNS[:numLimit] {
                nonRepeatedList = append(nonRepeatedList, a+b)
            }
        }
        // word + symbol
        for _, a := range basicPattern {
            for _, b := range SYMBOLIC_PATTERNS[:symLimit] {
                nonRepeatedList = append(nonRepeatedList, a+b)
            }
        }
        // word + number + symbol
        for _, a := range basicPattern {
            for _, b := range NUMERIC_PATTERNS[:numLimit] {
                for _, c := range SYMBOLIC_PATTERNS[:symLimit] {
                    nonRepeatedList = append(nonRepeatedList, a+b+c)
                }
            }
        }
        // word + symbol + number
        for _, a := range basicPattern {
            for _, b := range SYMBOLIC_PATTERNS[:symLimit] {
                for _, c := range NUMERIC_PATTERNS[:numLimit] {
                    nonRepeatedList = append(nonRepeatedList, a+b+c)
                }
            }
        }
        // symbol + word
        for _, a := range SYMBOLIC_PATTERNS[:symLimit] {
            for _, b := range basicPattern {
                nonRepeatedList = append(nonRepeatedList, a+b)
            }
        }
        // number + word
        for _, a := range NUMERIC_PATTERNS[:numLimit] {
            for _, b := range basicPattern {
                nonRepeatedList = append(nonRepeatedList, a+b)
            }
        }
        // symbol + word + number
        for _, a := range SYMBOLIC_PATTERNS[:symLimit] {
            for _, b := range basicPattern {
                for _, c := range NUMERIC_PATTERNS[:numLimit] {
                    nonRepeatedList = append(nonRepeatedList, a+b+c)
                }
            }
        }
        // number + word + symbol
        for _, a := range NUMERIC_PATTERNS[:numLimit] {
            for _, b := range basicPattern {
                for _, c := range SYMBOLIC_PATTERNS[:symLimit] {
                    nonRepeatedList = append(nonRepeatedList, a+b+c)
                }
            }
        }
		// word + suffix + symbol
		for _, a := range basicPattern {
			for _, b := range BASIC_SUFIXS[:limitSuf] {
				for _, c := range SYMBOLIC_PATTERNS[:symLimit] {
					nonRepeatedList = append(nonRepeatedList, a+b+c)
				}
			}
		}
    }

    // ✅ Dedupe al final (ORIGINAL)
    return uniquePreserve(nonRepeatedList)
}

func uniquePreserve(items []string) []string {
	seen := map[string]struct{}{}
	out := []string{}
	for _, it := range items {
		if _, ok := seen[it]; !ok {
			seen[it] = struct{}{}
			out = append(out, it)
		}
	}
	return out
}

// extract patterns from file (prefixes, suffixes, numbers, symbols)
func extractPatterns(fileInput string) ([]string, []string, []string, []string) {
	suf := []string{}
	pref := []string{}
	nums := []string{}
	syms := []string{}

	f, err := os.Open(fileInput)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening dictionary file: %v\n", err)
		return suf, pref, nums, syms
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	reWord := regexp.MustCompile(`(?:[^\W\d]|-){3,}`)
	reNums := regexp.MustCompile(`\d+`)
	reSyms := regexp.MustCompile(`[^\w\s]+`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		match := reWord.FindStringIndex(line)
		if match != nil {
			conceptStart := match[0]
			conceptEnd := match[1]
			prefix := strings.TrimSpace(line[:conceptStart])
			suffix := strings.TrimSpace(line[conceptEnd:])
			if prefix != "" {
				pref = append(pref, prefix)
			}
			if suffix != "" {
				suf = append(suf, suffix)
			}
			verbosePrint("[+] Prefix and Sufix successfully extracted")
		}
		foundNumbers := reNums.FindAllString(line, -1)
		nums = append(nums, foundNumbers...)
		verbosePrint("[+] Numeric patterns successfully extracted")
		foundSymbols := reSyms.FindAllString(line, -1)
		syms = append(syms, foundSymbols...)
		verbosePrint("[+] Symbol patterns successfully extracted")
	}
	// count and keep items that appear >=2
	countAndFilter := func(arr []string) []string {
		cnt := map[string]int{}
		for _, a := range arr {
			cnt[a]++
		}
		type kv struct {
			Key string
			V   int
		}
		list := []kv{}
		for k, v := range cnt {
			if v >= 2 {
				list = append(list, kv{k, v})
			}
		}
		sort.Slice(list, func(i, j int) bool {
			return list[i].V > list[j].V
		})
		out := []string{}
		for _, it := range list {
			out = append(out, it.Key)
		}
		return out
	}
	return countAndFilter(suf), countAndFilter(pref), countAndFilter(nums), countAndFilter(syms)
}

func main() {
	// parse args
	users := flag.String("u", "", "File with usernames, or usernames list: \"jony random,fahim jordan,...\"")
	password := flag.String("p", "", "file with words to 'passworize', or list like: \"ibis,megacorp,...\"")
	jn := flag.String("jn", "", "Use only the fasttext module to see nearest neighbours from words")
	light := flag.Bool("l", false, "Light mode, for small list (passwd mode).")
	full := flag.Bool("f", false, "Full mode. Warning, the output could be very heavy (passwd mode).")
	noVerbose := flag.Bool("nv", false, "Remove any output except the dictionary itself (Errors will be shown anyway).")
	dictionary := flag.String("d", "", "Extract patterns from your an specific dictionary.")
	output := flag.String("o", "", "Dictionary will be stored in this file.")
	ml := flag.String("ml", "", "Use a trained machine learning model to include neighbors of your original words.")
	n := flag.Int("n", 0, "Ammount of neighbors maximum for each word (20 by Default).")

	flag.Parse()

	if *light && *full {
		fmt.Fprintln(os.Stderr, "[!] Light mode and Full mode can not be at same time. Exiting...")
		os.Exit(1)
	}
	if *noVerbose {
		VERBOSE = false
	}
	if *light {
		LIGHT_MODE = true
		verbosePrint("[+] Light mode enabled.")
	}
	if *full {
		FULL_MODE = true
		verbosePrint("[+] Full mode enabled.")
		verbosePrint("[!] Warning, the output could be very heavy.")
	}
	outputFileName := ""
	if *output != "" {
		OUTPUT_FILE_BOOLEAN = true
		outputFileName = *output
		verbosePrint("[i] Dictionary will be stored in this file -> " + outputFileName)
	}

	if *dictionary != "" {
		if isAValidFile(*dictionary) {
			verbosePrint("[+] Using " + *dictionary + " as a dictionary.")
		} else {
			fmt.Fprintln(os.Stderr, "[!] File introduced as dictionary is not valid. Exiting...")
			os.Exit(1)
		}
		verbosePrint("[+] Extracting patterns... this will take a minut.")
		suf, pref, nums, syms := extractPatterns(*dictionary)
		if len(suf) > 0 {
			BASIC_SUFIXS = suf
		}
		if len(pref) > 0 {
			BASIC_PREFIXS = pref
		}
		if len(nums) > 0 {
			NUMERIC_PATTERNS = nums
		}
		if len(syms) > 0 {
			SYMBOLIC_PATTERNS = syms
		}
		verbosePrint("[+] Patterns successfuly extracted, creating dictionary...")
	}

	if *users != "" {
		if *ml != "" {
			fmt.Fprintln(os.Stderr, "[-] Sorry, machine-learning-model option is not compatible with USERS mode, is only for PASSWORD mode.")
			os.Exit(1)
		}
		if strings.TrimSpace(*users) == "" {
			fmt.Fprintln(os.Stderr, "Error: This argument can not be empty")
			os.Exit(1)
		}
		verbosePrint("[+] USER mode selected.")
		if detectIfFileOrNot(*users) {
			processFileUser(*users, outputFileName)
			os.Exit(0)
		} else {
			processInputUser(*users, outputFileName)
			os.Exit(0)
		}
	} else if *password != "" {
		if strings.TrimSpace(*password) == "" {
			fmt.Fprintln(os.Stderr, "Error: This argument can not be empty")
			os.Exit(1)
		}
		inputList := normalizeList(*password)
		if *ml != "" {
			if !detectIfFileOrNot(*ml) {
				fmt.Fprintln(os.Stderr, "[-] Your -ml <input> is not even a file...  ¬¬ , set a file here.")
				os.Exit(1)
			}
			if *n != 0 {
				NEIGHBORS_AMMOUNT = *n
			}
			ram, err := getTotalRAM()
			if err == nil {
				if ram < 11 {
					fmt.Fprintf(os.Stderr, "[-] Insufficient RAM for ML mode. Detected -> %.2fGB\n", ram)
					fmt.Fprintln(os.Stderr, "[-] You need 11 GB RAM available at least to run fasttext models.")
					os.Exit(1)
				}
			}
			mlList := mlProcessPwd(inputList, *ml, NEIGHBORS_AMMOUNT)
			processPasswd(mlList, outputFileName)
			os.Exit(0)
		}
		verbosePrint("[+] PASSWORD mode selected.")
		processPasswd(inputList, outputFileName)
		os.Exit(0)
	}

	if *jn != "" {
		if *ml == "" {
			fmt.Fprintln(os.Stderr, "[-] Sorry, you need to provide the Machine learning model (use -ml flag)")
			os.Exit(1)
		}
		if *n == 0 {
			fmt.Fprintln(os.Stderr, "[-] Sorry, you need to provide how many neighbours you are looking for (use -n flag)")
			os.Exit(1)
		}
		inputList := normalizeList(*jn)
		if *ml != "" {
			if !detectIfFileOrNot(*ml) {
				fmt.Fprintln(os.Stderr, "[-] Your -ml <input> is not even a file...  ¬¬ , set a file here.")
				os.Exit(1)
			}
			if *n != 0 {
				NEIGHBORS_AMMOUNT = *n
			}
			mlList := mlProcessPwd(inputList, *ml, NEIGHBORS_AMMOUNT)
			fmt.Println(mlList)
			os.Exit(0)
		}
	}

	// default: print banner if no args
	if len(os.Args) == 1 {
		printBanner()
	}
}
