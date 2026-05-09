@echo off
setlocal enabledelayedexpansion

set "BASE_URL=http://localhost:8080"
set "USER=adminins"
set "PASS=12345678910"
set "EMAIL=admin2@gmail.com"

echo ========================================
echo 1. DANG KY TAI KHOAN...
curl -s -X POST "%BASE_URL%/auth/register" ^
     -H "Content-Type: application/json" ^
     -d "{\"username\": \"%USER%\", \"email\": \"%EMAIL%\", \"password\": \"%PASS%\", \"RePassword\": \"%PASS%\"}"
echo.

echo ========================================
echo 2. DANG NHAP...
curl -s -X POST "%BASE_URL%/auth/login" ^
     -H "Content-Type: application/json" ^
     -d "{\"username\": \"%USER%\", \"password\": \"%PASS%\"}" > login_res.txt

for /f "usebackq delims=" %%a in (`powershell -Command "$j = Get-Content login_res.txt | ConvertFrom-Json; Write-Output $j.token"`) do set "TOKEN=%%a"
if exist login_res.txt del login_res.txt

if "!TOKEN!"=="" (
    echo [LOI] Khong lay duoc Token.
    pause & exit /b
)
echo Token: !TOKEN:~0,25!...
echo.

set "MANGA[1]=Naruto|Kishimoto Masashi|Action,Adventure,Shounen|completed|700|A young ninja seeks to become Hokage"
set "MANGA[2]=One Piece|Eiichiro Oda|Action,Adventure,Shounen|ongoing|1100|A pirate searches for the ultimate treasure"
set "MANGA[3]=Bleach|Tite Kubo|Action,Adventure,Shounen|completed|686|A teen gains soul reaper powers"
set "MANGA[4]=Attack on Titan|Hajime Isayama|Action,Drama,Seinen|completed|139|Humanity fights giants behind walls"
set "MANGA[5]=Death Note|Tsugumi Ohba|Mystery,Psychological,Thriller|completed|108|A student uses a supernatural notebook"
set "MANGA[6]=Fullmetal Alchemist|Hiromu Arakawa|Action,Adventure,Fantasy|completed|108|Two brothers seek the Philosophers Stone"
set "MANGA[7]=Dragon Ball|Akira Toriyama|Action,Adventure,Shounen|completed|519|A warrior trains to become the strongest"
set "MANGA[8]=Demon Slayer|Koyoharu Gotouge|Action,Adventure,Shounen|completed|205|A boy hunts demons to save his sister"
set "MANGA[9]=My Hero Academia|Kohei Horikoshi|Action,Adventure,Shounen|ongoing|420|A powerless boy trains to become a hero"
set "MANGA[10]=Tokyo Ghoul|Sui Ishida|Action,Horror,Seinen|completed|179|A student becomes half ghoul after an attack"
set "MANGA[11]=Sword Art Online|Reki Kawahara|Action,Adventure,Fantasy|ongoing|26|Players trapped in a virtual reality game"
set "MANGA[12]=Hunter x Hunter|Yoshihiro Togashi|Action,Adventure,Shounen|ongoing|400|A boy seeks his missing hunter father"
set "MANGA[13]=Fairy Tail|Hiro Mashima|Action,Adventure,Fantasy|completed|545|A wizard joins a powerful guild of mages"
set "MANGA[14]=Black Clover|Yuki Tabata|Action,Adventure,Shounen|ongoing|370|A boy with no magic dreams of becoming wizard king"
set "MANGA[15]=Jujutsu Kaisen|Gege Akutami|Action,Horror,Shounen|ongoing|260|A student fights cursed spirits after swallowing a finger"
set "MANGA[16]=Vinland Saga|Makoto Yukimura|Action,Adventure,Historical|ongoing|200|A Viking seeks revenge for his fathers death"
set "MANGA[17]=Berserk|Kentaro Miura|Action,Adventure,Horror|ongoing|374|A lone swordsman fights demons in a dark world"
set "MANGA[18]=Vagabond|Takehiko Inoue|Action,Historical,Seinen|ongoing|327|The story of legendary samurai Miyamoto Musashi"
set "MANGA[19]=Slam Dunk|Takehiko Inoue|Sports,Comedy,Shounen|completed|276|A delinquent joins basketball to impress a girl"
set "MANGA[20]=Haikyuu|Haruichi Furudate|Sports,Comedy,Shounen|completed|402|A short boy dreams of becoming a volleyball ace"
set "MANGA[21]=Kuroko no Basket|Tadatoshi Fujimaki|Sports,Shounen|completed|276|A phantom player joins a high school basketball team"
set "MANGA[22]=Captain Tsubasa|Yoichi Takahashi|Sports,Shounen|ongoing|130|A boy dreams of becoming the worlds best soccer player"
set "MANGA[23]=Eyeshield 21|Riichiro Inagaki|Sports,Comedy,Shounen|completed|333|A timid boy becomes a star American football player"
set "MANGA[24]=Hajime no Ippo|George Morikawa|Sports,Shounen|ongoing|1400|A bullied teen becomes a professional boxer"
set "MANGA[25]=Yowamushi Pedal|Wataru Watanabe|Sports,Shounen|ongoing|740|An anime fan discovers his talent for cycling"
set "MANGA[26]=Shingeki no Kyojin|Hajime Isayama|Action,Drama,Seinen|completed|139|Soldiers fight titans threatening humanitys survival"
set "MANGA[27]=Re Zero|Tappei Nagatsuki|Fantasy,Isekai,Drama|ongoing|90|A boy transported to another world dies and respawns"
set "MANGA[28]=Overlord|Kugane Maruyama|Fantasy,Isekai,Action|ongoing|80|A player trapped in a game as an undead overlord"
set "MANGA[29]=Slime Isekai|Fuse|Fantasy,Isekai,Comedy|ongoing|120|A man reincarnates as a powerful slime"
set "MANGA[30]=Konosuba|Natsume Akatsuki|Fantasy,Isekai,Comedy|completed|20|A boy and useless goddess adventure in another world"
set "MANGA[31]=SAO Progressive|Reki Kawahara|Action,Fantasy,Romance|ongoing|12|Detailed retelling of the Aincrad arc"
set "MANGA[32]=No Game No Life|Yuu Kamiya|Fantasy,Isekai,Comedy|ongoing|10|Genius siblings transported to a world of games"
set "MANGA[33]=Shield Hero|Aneko Yusagi|Fantasy,Isekai,Action|ongoing|90|A falsely accused hero rebuilds his strength"
set "MANGA[34]=Mushoku Tensei|Rifujin na Magonote|Fantasy,Isekai,Drama|ongoing|90|A man reincarnates and commits to a new life"
set "MANGA[35]=Made in Abyss|Akihito Tsukushi|Adventure,Fantasy,Horror|ongoing|67|A girl dives into a mysterious abyss to find her mother"
set "MANGA[36]=Vinland Saga Prologue|Makoto Yukimura|Action,Historical|completed|54|Early adventures of Thorfinn in Europe"
set "MANGA[37]=Pluto|Naoki Urasawa|Mystery,Sci-fi,Drama|completed|65|A robot detective investigates a series of murders"
set "MANGA[38]=Monster|Naoki Urasawa|Mystery,Psychological,Thriller|completed|162|A doctor hunts a serial killer he once saved"
set "MANGA[39]=20th Century Boys|Naoki Urasawa|Mystery,Sci-fi,Drama|completed|249|Friends uncover a global conspiracy from their childhood"
set "MANGA[40]=Gantz|Hiroya Oku|Action,Sci-fi,Horror|completed|383|Dead people forced to fight aliens for survival"
set "MANGA[41]=Elfen Lied|Lynn Okamoto|Action,Drama,Horror|completed|107|A mutant girl escapes a lab and causes chaos"
set "MANGA[42]=Claymore|Norihiro Yagi|Action,Adventure,Fantasy|completed|155|Female warriors fight demons in a medieval world"
set "MANGA[43]=Inuyasha|Rumiko Takahashi|Action,Adventure,Romance|completed|558|A girl travels feudal Japan with a half demon"
set "MANGA[44]=Ranma Half|Rumiko Takahashi|Action,Comedy,Romance|completed|407|A boy transforms into a girl when splashed with cold water"
set "MANGA[45]=Yu Yu Hakusho|Yoshihiro Togashi|Action,Adventure,Shounen|completed|175|A delinquent becomes a spirit detective"
set "MANGA[46]=Saint Seiya|Masami Kurumada|Action,Adventure,Shounen|completed|328|Young warriors protect the goddess Athena"
set "MANGA[47]=Rurouni Kenshin|Nobuhiro Watsuki|Action,Historical,Romance|completed|255|A former assassin seeks a peaceful life in Meiji Japan"
set "MANGA[48]=Samurai X|Nobuhiro Watsuki|Action,Historical|completed|96|Prequel story of Kenshin as a ruthless assassin"
set "MANGA[49]=Trigun|Yasuhiro Nightow|Action,Sci-fi,Comedy|completed|84|A pacifist gunman is hunted across a desert planet"
set "MANGA[50]=Cowboy Bebop|Yutaka Nanten|Action,Sci-fi,Adventure|completed|3|Bounty hunters travel the solar system in 2071"
set "MANGA[51]=Seraph of the End|Takaya Kagami|Action,Horror,Shounen|ongoing|130|Humans fight vampires who took over the world"
set "MANGA[52]=Blue Exorcist|Kazue Kato|Action,Fantasy,Shounen|ongoing|140|Son of Satan trains to become an exorcist"
set "MANGA[53]=Noragami|Adachitoka|Action,Fantasy,Shounen|ongoing|105|A minor god seeks fame and a proper shrine"
set "MANGA[54]=Soul Eater|Atsushi Ohkubo|Action,Fantasy,Shounen|completed|113|Students at a school for weapon meisters fight evil"
set "MANGA[55]=Fire Force|Atsushi Ohkubo|Action,Fantasy,Shounen|completed|304|Firefighters battle humans who spontaneously combust"
set "MANGA[56]=Toriko|Mitsutoshi Shimabukuro|Action,Adventure,Shounen|completed|396|A hunter seeks the ultimate ingredients for a meal"
set "MANGA[57]=Beelzebub|Ryuhei Tamura|Action,Comedy,Shounen|completed|240|A delinquent is chosen to raise the demon kings baby"
set "MANGA[58]=Katekyo Hitman Reborn|Akira Amano|Action,Comedy,Shounen|completed|409|A weak student trains to become a mafia boss"
set "MANGA[59]=Medaka Box|NisiOisin|Action,Comedy,Shounen|completed|192|A perfect student council president solves school problems"
set "MANGA[60]=Nisekoi|Naoshi Komi|Comedy,Romance,Shounen|completed|229|A gangster son and a yakuza daughter fake a relationship"
set "MANGA[61]=Oregairu|Wataru Watari|Comedy,Romance,School|completed|14|A cynical loner joins a service club at school"
set "MANGA[62]=Clannad|Key|Drama,Romance,School|completed|4|A delinquent meets a girl repeating her school year"
set "MANGA[63]=Kanon|Key|Drama,Romance,School|completed|5|A boy returns to a snowy town and meets mysterious girls"
set "MANGA[64]=Angel Beats|Jun Maeda|Drama,Comedy,School|completed|3|Students in the afterlife fight against God"
set "MANGA[65]=Charlotte|Jun Maeda|Drama,Sci-fi,School|completed|7|Teenagers with special abilities attend a school together"
set "MANGA[66]=A Silent Voice|Yoshitoki Oima|Drama,Romance,School|completed|62|A former bully seeks redemption with a deaf girl"
set "MANGA[67]=Your Lie in April|Naoshi Arakawa|Drama,Romance,Music|completed|44|A pianist meets a violinist who changes his world"
set "MANGA[68]=Nana|Ai Yazawa|Drama,Romance,Music|ongoing|84|Two women named Nana share an apartment in Tokyo"
set "MANGA[69]=Paradise Kiss|Ai Yazawa|Drama,Romance,Fashion|completed|40|A student falls in love with a fashion designer"
set "MANGA[70]=Fruits Basket|Natsuki Takaya|Drama,Romance,Shoujo|completed|136|A girl discovers a family cursed by the zodiac"
set "MANGA[71]=Cardcaptor Sakura|CLAMP|Adventure,Fantasy,Shoujo|completed|50|A girl captures magical cards scattered across her town"
set "MANGA[72]=Sailor Moon|Naoko Takeuchi|Action,Fantasy,Shoujo|completed|60|A teenage girl transforms into a guardian of love"
set "MANGA[73]=Utena|Chiho Saito|Drama,Fantasy,Shoujo|completed|39|A girl duels to protect a mysterious rose bride"
set "MANGA[74]=Ouran Host Club|Bisco Hatori|Comedy,Romance,Shoujo|completed|87|A girl joins a host club to repay a debt"
set "MANGA[75]=Skip Beat|Yoshiki Nakamura|Comedy,Drama,Shoujo|ongoing|300|A girl enters showbiz to take revenge on her ex"
set "MANGA[76]=Kamisama Kiss|Julietta Suzuki|Comedy,Fantasy,Shoujo|completed|148|A homeless girl becomes a land god"
set "MANGA[77]=Maid-sama|Hiro Fujiwara|Comedy,Romance,Shoujo|completed|85|A strict student council president hides her maid cafe job"
set "MANGA[78]=Ao Haru Ride|Io Sakisaka|Drama,Romance,Shoujo|completed|49|A girl reunites with her first love in high school"
set "MANGA[79]=Strobe Edge|Io Sakisaka|Drama,Romance,Shoujo|completed|10|A girl falls for the most popular boy in school"
set "MANGA[80]=Bokura ga Ita|Yuuki Obata|Drama,Romance,Shoujo|completed|16|A girl falls for a boy with a complicated past"
set "MANGA[81]=Mushishi|Yuki Urushibara|Adventure,Mystery,Seinen|completed|50|A traveler helps people afflicted by mysterious creatures"
set "MANGA[82]=Dungeon Meshi|Ryoko Kui|Adventure,Fantasy,Comedy|completed|97|Adventurers cook and eat monsters in a dungeon"
set "MANGA[83]=Dorohedoro|Q Hayashida|Action,Fantasy,Comedy|completed|167|A man with a reptile head searches for his true identity"
set "MANGA[84]=Biomega|Tsutomu Nihei|Action,Sci-fi,Horror|completed|42|A synthetic human rides a motorcycle through a dying world"
set "MANGA[85]=Blame|Tsutomu Nihei|Action,Sci-fi,Horror|completed|66|A wanderer searches for humans in a massive megastructure"
set "MANGA[86]=Holyland|Kouji Mori|Action,Drama,Seinen|completed|182|A bullied teen finds peace through street fighting"
set "MANGA[87]=Homunculus|Hideo Yamamoto|Drama,Horror,Seinen|completed|166|A man gains the ability to see peoples traumas"
set "MANGA[88]=Oyasumi Punpun|Inio Asano|Drama,Psychological,Seinen|completed|147|A surreal coming of age story of a troubled boy"
set "MANGA[89]=Solanin|Inio Asano|Drama,Romance,Seinen|completed|40|Young adults struggle with life after college"
set "MANGA[90]=I Am a Hero|Kengo Hanazawa|Action,Horror,Seinen|completed|22|A manga artist survives a zombie apocalypse in Japan"
set "MANGA[91]=Deadman Wonderland|Jinsei Kataoka|Action,Horror,Sci-fi|completed|58|A boy is sent to a prison amusement park after being framed"
set "MANGA[92]=Mirai Nikki|Sakae Esuno|Action,Mystery,Thriller|completed|59|A boy receives a diary that predicts the future"
set "MANGA[93]=Highschool of the Dead|Daisuke Sato|Action,Horror,Drama|completed|29|Students fight to survive a sudden zombie outbreak"
set "MANGA[94]=Terra Formars|Yu Sasuga|Action,Sci-fi,Horror|ongoing|230|Humans fight evolved cockroaches on Mars"
set "MANGA[95]=Akame ga Kill|Takahiro|Action,Drama,Fantasy|completed|78|A naive boy joins an assassin group fighting corruption"
set "MANGA[96]=Seven Deadly Sins|Nakaba Suzuki|Action,Adventure,Fantasy|completed|346|A princess seeks legendary knights to save her kingdom"
set "MANGA[97]=Magi|Shinobu Ohtaka|Action,Adventure,Fantasy|completed|372|A young mage explores dungeons to find his destiny"
set "MANGA[98]=Radiant|Tony Valente|Action,Adventure,Fantasy|ongoing|21|A young sorcerer hunts the source of all monsters"
set "MANGA[99]=Witch Hat Atelier|Kamome Shirahama|Adventure,Fantasy,Drama|ongoing|14|A girl discovers magic is learned not born"
set "MANGA[100]=Golden Kamuy|Satoru Noda|Action,Adventure,Historical|completed|314|A soldier hunts for hidden Ainu gold in Hokkaido"

set "COUNT=100"

echo ========================================
echo 3. THEM 100 MANGA...
echo.

for /l %%i in (1,1,%COUNT%) do (
    for /f "tokens=1-6 delims=|" %%a in ("!MANGA[%%i]!") do (
        set "TITLE=%%a"
        set "AUTHOR=%%b"
        set "GENRES=%%c"
        set "STATUS=%%d"
        set "CHAPTERS=%%e"
        set "DESC=%%f"

        for /f "usebackq delims=" %%g in (`powershell -Command "$g = '!GENRES!' -split ','; $arr = $g | ForEach-Object { '\"' + $_.Trim() + '\"' }; '[' + ($arr -join ',') + ']'"`) do set "GENRES_JSON=%%g"

        echo [%%i/100] Adding: !TITLE!
        curl -s -X POST "%BASE_URL%/manga" ^
             -H "Content-Type: application/json" ^
             -H "Authorization: Bearer !TOKEN!" ^
             -d "{\"title\": \"!TITLE!\", \"author\": \"!AUTHOR!\", \"genres\": !GENRES_JSON!, \"status\": \"!STATUS!\", \"total_chapters\": !CHAPTERS!, \"description\": \"!DESC!\"}"
        echo.
    )
)

echo.
echo ========================================
echo DANG XUAT...
curl -s -X POST "%BASE_URL%/auth/logout" ^
     -H "Authorization: Bearer !TOKEN!"

echo.
echo DA HOAN THANH 100 MANGA!
pause