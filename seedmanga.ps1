$BASE_URL = "http://localhost:8080"
$USER = "adminins"
$PASS = "12345678910"
$EMAIL = "admin2@gmail.com"

# 1. Dang ky
Write-Host "========================================"
Write-Host "1. DANG KY TAI KHOAN..."
try {
    Invoke-RestMethod -Method POST "$BASE_URL/auth/register" `
        -ContentType "application/json" `
        -Body (@{ username=$USER; email=$EMAIL; password=$PASS; RePassword=$PASS } | ConvertTo-Json)
} catch { Write-Host "Da ton tai hoac loi dang ky" }

# 2. Dang nhap
Write-Host "========================================"
Write-Host "2. DANG NHAP..."
$loginRes = Invoke-RestMethod -Method POST "$BASE_URL/auth/login" `
    -ContentType "application/json" `
    -Body (@{ username=$USER; password=$PASS } | ConvertTo-Json)

$TOKEN = $loginRes.token
if (-not $TOKEN) {
    Write-Host "[LOI] Khong lay duoc Token."
    exit 1
}
Write-Host "Dang nhap thanh cong! Token: $($TOKEN.Substring(0,25))..."

$headers = @{ Authorization = "Bearer $TOKEN" }

# 3. Danh sach manga
$mangas = @(
    @{ title="Naruto"; author="Kishimoto Masashi"; genres=@("Action","Adventure","Shounen"); status="completed"; total_chapters=700; description="A young ninja seeks to become Hokage" },
    @{ title="One Piece"; author="Eiichiro Oda"; genres=@("Action","Adventure","Shounen"); status="ongoing"; total_chapters=1100; description="A pirate searches for the ultimate treasure" },
    @{ title="Bleach"; author="Tite Kubo"; genres=@("Action","Adventure","Shounen"); status="completed"; total_chapters=686; description="A teen gains soul reaper powers" },
    @{ title="Attack on Titan"; author="Hajime Isayama"; genres=@("Action","Drama","Seinen"); status="completed"; total_chapters=139; description="Humanity fights giants behind walls" },
    @{ title="Death Note"; author="Tsugumi Ohba"; genres=@("Mystery","Psychological","Thriller"); status="completed"; total_chapters=108; description="A student uses a supernatural notebook to kill criminals" },
    @{ title="Fullmetal Alchemist"; author="Hiromu Arakawa"; genres=@("Action","Adventure","Fantasy"); status="completed"; total_chapters=108; description="Two brothers seek the Philosophers Stone" },
    @{ title="Dragon Ball"; author="Akira Toriyama"; genres=@("Action","Adventure","Shounen"); status="completed"; total_chapters=519; description="A warrior trains to become the strongest" },
    @{ title="Demon Slayer"; author="Koyoharu Gotouge"; genres=@("Action","Adventure","Shounen"); status="completed"; total_chapters=205; description="A boy hunts demons to save his sister" },
    @{ title="My Hero Academia"; author="Kohei Horikoshi"; genres=@("Action","Adventure","Shounen"); status="ongoing"; total_chapters=420; description="A powerless boy trains to become a hero" },
    @{ title="Tokyo Ghoul"; author="Sui Ishida"; genres=@("Action","Horror","Seinen"); status="completed"; total_chapters=179; description="A student becomes half ghoul after an attack" },
    @{ title="Sword Art Online"; author="Reki Kawahara"; genres=@("Action","Adventure","Fantasy"); status="ongoing"; total_chapters=26; description="Players trapped in a virtual reality game" },
    @{ title="Hunter x Hunter"; author="Yoshihiro Togashi"; genres=@("Action","Adventure","Shounen"); status="ongoing"; total_chapters=400; description="A boy seeks his missing hunter father" },
    @{ title="Fairy Tail"; author="Hiro Mashima"; genres=@("Action","Adventure","Fantasy"); status="completed"; total_chapters=545; description="A wizard joins a powerful guild of mages" },
    @{ title="Black Clover"; author="Yuki Tabata"; genres=@("Action","Adventure","Shounen"); status="ongoing"; total_chapters=370; description="A boy with no magic dreams of becoming wizard king" },
    @{ title="Jujutsu Kaisen"; author="Gege Akutami"; genres=@("Action","Horror","Shounen"); status="ongoing"; total_chapters=260; description="A student fights cursed spirits after swallowing a finger" },
    @{ title="Vinland Saga"; author="Makoto Yukimura"; genres=@("Action","Adventure","Historical"); status="ongoing"; total_chapters=200; description="A Viking seeks revenge for his fathers death" },
    @{ title="Berserk"; author="Kentaro Miura"; genres=@("Action","Adventure","Horror"); status="ongoing"; total_chapters=374; description="A lone swordsman fights demons in a dark world" },
    @{ title="Vagabond"; author="Takehiko Inoue"; genres=@("Action","Historical","Seinen"); status="ongoing"; total_chapters=327; description="The story of legendary samurai Miyamoto Musashi" },
    @{ title="Slam Dunk"; author="Takehiko Inoue"; genres=@("Sports","Comedy","Shounen"); status="completed"; total_chapters=276; description="A delinquent joins basketball to impress a girl" },
    @{ title="Haikyuu"; author="Haruichi Furudate"; genres=@("Sports","Comedy","Shounen"); status="completed"; total_chapters=402; description="A short boy dreams of becoming a volleyball ace" },
    @{ title="Kuroko no Basket"; author="Tadatoshi Fujimaki"; genres=@("Sports","Shounen"); status="completed"; total_chapters=276; description="A phantom player joins a high school basketball team" },
    @{ title="Captain Tsubasa"; author="Yoichi Takahashi"; genres=@("Sports","Shounen"); status="ongoing"; total_chapters=130; description="A boy dreams of becoming the worlds best soccer player" },
    @{ title="Eyeshield 21"; author="Riichiro Inagaki"; genres=@("Sports","Comedy","Shounen"); status="completed"; total_chapters=333; description="A timid boy becomes a star American football player" },
    @{ title="Hajime no Ippo"; author="George Morikawa"; genres=@("Sports","Shounen"); status="ongoing"; total_chapters=1400; description="A bullied teen becomes a professional boxer" },
    @{ title="Yowamushi Pedal"; author="Wataru Watanabe"; genres=@("Sports","Shounen"); status="ongoing"; total_chapters=740; description="An anime fan discovers his talent for cycling" },
    @{ title="Re Zero"; author="Tappei Nagatsuki"; genres=@("Fantasy","Isekai","Drama"); status="ongoing"; total_chapters=90; description="A boy transported to another world dies and respawns" },
    @{ title="Overlord"; author="Kugane Maruyama"; genres=@("Fantasy","Isekai","Action"); status="ongoing"; total_chapters=80; description="A player trapped in a game as an undead overlord" },
    @{ title="Slime Isekai"; author="Fuse"; genres=@("Fantasy","Isekai","Comedy"); status="ongoing"; total_chapters=120; description="A man reincarnates as a powerful slime" },
    @{ title="Konosuba"; author="Natsume Akatsuki"; genres=@("Fantasy","Isekai","Comedy"); status="completed"; total_chapters=20; description="A boy and useless goddess adventure in another world" },
    @{ title="No Game No Life"; author="Yuu Kamiya"; genres=@("Fantasy","Isekai","Comedy"); status="ongoing"; total_chapters=10; description="Genius siblings transported to a world of games" },
    @{ title="Shield Hero"; author="Aneko Yusagi"; genres=@("Fantasy","Isekai","Action"); status="ongoing"; total_chapters=90; description="A falsely accused hero rebuilds his strength" },
    @{ title="Mushoku Tensei"; author="Rifujin na Magonote"; genres=@("Fantasy","Isekai","Drama"); status="ongoing"; total_chapters=90; description="A man reincarnates and commits to a new life" },
    @{ title="Made in Abyss"; author="Akihito Tsukushi"; genres=@("Adventure","Fantasy","Horror"); status="ongoing"; total_chapters=67; description="A girl dives into a mysterious abyss to find her mother" },
    @{ title="Pluto"; author="Naoki Urasawa"; genres=@("Mystery","Sci-fi","Drama"); status="completed"; total_chapters=65; description="A robot detective investigates a series of murders" },
    @{ title="Monster"; author="Naoki Urasawa"; genres=@("Mystery","Psychological","Thriller"); status="completed"; total_chapters=162; description="A doctor hunts a serial killer he once saved" },
    @{ title="20th Century Boys"; author="Naoki Urasawa"; genres=@("Mystery","Sci-fi","Drama"); status="completed"; total_chapters=249; description="Friends uncover a global conspiracy from their childhood" },
    @{ title="Gantz"; author="Hiroya Oku"; genres=@("Action","Sci-fi","Horror"); status="completed"; total_chapters=383; description="Dead people forced to fight aliens for survival" },
    @{ title="Elfen Lied"; author="Lynn Okamoto"; genres=@("Action","Drama","Horror"); status="completed"; total_chapters=107; description="A mutant girl escapes a lab and causes chaos" },
    @{ title="Claymore"; author="Norihiro Yagi"; genres=@("Action","Adventure","Fantasy"); status="completed"; total_chapters=155; description="Female warriors fight demons in a medieval world" },
    @{ title="Inuyasha"; author="Rumiko Takahashi"; genres=@("Action","Adventure","Romance"); status="completed"; total_chapters=558; description="A girl travels feudal Japan with a half demon" },
    @{ title="Ranma Half"; author="Rumiko Takahashi"; genres=@("Action","Comedy","Romance"); status="completed"; total_chapters=407; description="A boy transforms into a girl when splashed with cold water" },
    @{ title="Yu Yu Hakusho"; author="Yoshihiro Togashi"; genres=@("Action","Adventure","Shounen"); status="completed"; total_chapters=175; description="A delinquent becomes a spirit detective" },
    @{ title="Saint Seiya"; author="Masami Kurumada"; genres=@("Action","Adventure","Shounen"); status="completed"; total_chapters=328; description="Young warriors protect the goddess Athena" },
    @{ title="Rurouni Kenshin"; author="Nobuhiro Watsuki"; genres=@("Action","Historical","Romance"); status="completed"; total_chapters=255; description="A former assassin seeks a peaceful life in Meiji Japan" },
    @{ title="Samurai X"; author="Nobuhiro Watsuki"; genres=@("Action","Historical"); status="completed"; total_chapters=96; description="Prequel story of Kenshin as a ruthless assassin" },
    @{ title="Trigun"; author="Yasuhiro Nightow"; genres=@("Action","Sci-fi","Comedy"); status="completed"; total_chapters=84; description="A pacifist gunman is hunted across a desert planet" },
    @{ title="Seraph of the End"; author="Takaya Kagami"; genres=@("Action","Horror","Shounen"); status="ongoing"; total_chapters=130; description="Humans fight vampires who took over the world" },
    @{ title="Blue Exorcist"; author="Kazue Kato"; genres=@("Action","Fantasy","Shounen"); status="ongoing"; total_chapters=140; description="Son of Satan trains to become an exorcist" },
    @{ title="Noragami"; author="Adachitoka"; genres=@("Action","Fantasy","Shounen"); status="ongoing"; total_chapters=105; description="A minor god seeks fame and a proper shrine" },
    @{ title="Soul Eater"; author="Atsushi Ohkubo"; genres=@("Action","Fantasy","Shounen"); status="completed"; total_chapters=113; description="Students at a school for weapon meisters fight evil" },
    @{ title="Fire Force"; author="Atsushi Ohkubo"; genres=@("Action","Fantasy","Shounen"); status="completed"; total_chapters=304; description="Firefighters battle humans who spontaneously combust" },
    @{ title="Toriko"; author="Mitsutoshi Shimabukuro"; genres=@("Action","Adventure","Shounen"); status="completed"; total_chapters=396; description="A hunter seeks the ultimate ingredients for a meal" },
    @{ title="Beelzebub"; author="Ryuhei Tamura"; genres=@("Action","Comedy","Shounen"); status="completed"; total_chapters=240; description="A delinquent is chosen to raise the demon kings baby" },
    @{ title="Katekyo Hitman Reborn"; author="Akira Amano"; genres=@("Action","Comedy","Shounen"); status="completed"; total_chapters=409; description="A weak student trains to become a mafia boss" },
    @{ title="Nisekoi"; author="Naoshi Komi"; genres=@("Comedy","Romance","Shounen"); status="completed"; total_chapters=229; description="A gangster son and a yakuza daughter fake a relationship" },
    @{ title="Oregairu"; author="Wataru Watari"; genres=@("Comedy","Romance","School"); status="completed"; total_chapters=14; description="A cynical loner joins a service club at school" },
    @{ title="A Silent Voice"; author="Yoshitoki Oima"; genres=@("Drama","Romance","School"); status="completed"; total_chapters=62; description="A former bully seeks redemption with a deaf girl" },
    @{ title="Your Lie in April"; author="Naoshi Arakawa"; genres=@("Drama","Romance","Music"); status="completed"; total_chapters=44; description="A pianist meets a violinist who changes his world" },
    @{ title="Nana"; author="Ai Yazawa"; genres=@("Drama","Romance","Music"); status="ongoing"; total_chapters=84; description="Two women named Nana share an apartment in Tokyo" },
    @{ title="Fruits Basket"; author="Natsuki Takaya"; genres=@("Drama","Romance","Shoujo"); status="completed"; total_chapters=136; description="A girl discovers a family cursed by the zodiac" },
    @{ title="Cardcaptor Sakura"; author="CLAMP"; genres=@("Adventure","Fantasy","Shoujo"); status="completed"; total_chapters=50; description="A girl captures magical cards scattered across her town" },
    @{ title="Sailor Moon"; author="Naoko Takeuchi"; genres=@("Action","Fantasy","Shoujo"); status="completed"; total_chapters=60; description="A teenage girl transforms into a guardian of love" },
    @{ title="Ouran Host Club"; author="Bisco Hatori"; genres=@("Comedy","Romance","Shoujo"); status="completed"; total_chapters=87; description="A girl joins a host club to repay a debt" },
    @{ title="Skip Beat"; author="Yoshiki Nakamura"; genres=@("Comedy","Drama","Shoujo"); status="ongoing"; total_chapters=300; description="A girl enters showbiz to take revenge on her ex" },
    @{ title="Kamisama Kiss"; author="Julietta Suzuki"; genres=@("Comedy","Fantasy","Shoujo"); status="completed"; total_chapters=148; description="A homeless girl becomes a land god" },
    @{ title="Maid-sama"; author="Hiro Fujiwara"; genres=@("Comedy","Romance","Shoujo"); status="completed"; total_chapters=85; description="A strict student council president hides her maid cafe job" },
    @{ title="Ao Haru Ride"; author="Io Sakisaka"; genres=@("Drama","Romance","Shoujo"); status="completed"; total_chapters=49; description="A girl reunites with her first love in high school" },
    @{ title="Strobe Edge"; author="Io Sakisaka"; genres=@("Drama","Romance","Shoujo"); status="completed"; total_chapters=10; description="A girl falls for the most popular boy in school" },
    @{ title="Mushishi"; author="Yuki Urushibara"; genres=@("Adventure","Mystery","Seinen"); status="completed"; total_chapters=50; description="A traveler helps people afflicted by mysterious creatures" },
    @{ title="Dungeon Meshi"; author="Ryoko Kui"; genres=@("Adventure","Fantasy","Comedy"); status="completed"; total_chapters=97; description="Adventurers cook and eat monsters in a dungeon" },
    @{ title="Dorohedoro"; author="Q Hayashida"; genres=@("Action","Fantasy","Comedy"); status="completed"; total_chapters=167; description="A man with a reptile head searches for his true identity" },
    @{ title="Biomega"; author="Tsutomu Nihei"; genres=@("Action","Sci-fi","Horror"); status="completed"; total_chapters=42; description="A synthetic human rides a motorcycle through a dying world" },
    @{ title="Blame"; author="Tsutomu Nihei"; genres=@("Action","Sci-fi","Horror"); status="completed"; total_chapters=66; description="A wanderer searches for humans in a massive megastructure" },
    @{ title="Holyland"; author="Kouji Mori"; genres=@("Action","Drama","Seinen"); status="completed"; total_chapters=182; description="A bullied teen finds peace through street fighting" },
    @{ title="Homunculus"; author="Hideo Yamamoto"; genres=@("Drama","Horror","Seinen"); status="completed"; total_chapters=166; description="A man gains the ability to see peoples traumas" },
    @{ title="Oyasumi Punpun"; author="Inio Asano"; genres=@("Drama","Psychological","Seinen"); status="completed"; total_chapters=147; description="A surreal coming of age story of a troubled boy" },
    @{ title="Solanin"; author="Inio Asano"; genres=@("Drama","Romance","Seinen"); status="completed"; total_chapters=40; description="Young adults struggle with life after college" },
    @{ title="I Am a Hero"; author="Kengo Hanazawa"; genres=@("Action","Horror","Seinen"); status="completed"; total_chapters=22; description="A manga artist survives a zombie apocalypse in Japan" },
    @{ title="Deadman Wonderland"; author="Jinsei Kataoka"; genres=@("Action","Horror","Sci-fi"); status="completed"; total_chapters=58; description="A boy is sent to a prison amusement park after being framed" },
    @{ title="Mirai Nikki"; author="Sakae Esuno"; genres=@("Action","Mystery","Thriller"); status="completed"; total_chapters=59; description="A boy receives a diary that predicts the future" },
    @{ title="Highschool of the Dead"; author="Daisuke Sato"; genres=@("Action","Horror","Drama"); status="completed"; total_chapters=29; description="Students fight to survive a sudden zombie outbreak" },
    @{ title="Terra Formars"; author="Yu Sasuga"; genres=@("Action","Sci-fi","Horror"); status="ongoing"; total_chapters=230; description="Humans fight evolved cockroaches on Mars" },
    @{ title="Akame ga Kill"; author="Takahiro"; genres=@("Action","Drama","Fantasy"); status="completed"; total_chapters=78; description="A naive boy joins an assassin group fighting corruption" },
    @{ title="Seven Deadly Sins"; author="Nakaba Suzuki"; genres=@("Action","Adventure","Fantasy"); status="completed"; total_chapters=346; description="A princess seeks legendary knights to save her kingdom" },
    @{ title="Magi"; author="Shinobu Ohtaka"; genres=@("Action","Adventure","Fantasy"); status="completed"; total_chapters=372; description="A young mage explores dungeons to find his destiny" },
    @{ title="Radiant"; author="Tony Valente"; genres=@("Action","Adventure","Fantasy"); status="ongoing"; total_chapters=21; description="A young sorcerer hunts the source of all monsters" },
    @{ title="Witch Hat Atelier"; author="Kamome Shirahama"; genres=@("Adventure","Fantasy","Drama"); status="ongoing"; total_chapters=14; description="A girl discovers magic is learned not born" },
    @{ title="Golden Kamuy"; author="Satoru Noda"; genres=@("Action","Adventure","Historical"); status="completed"; total_chapters=314; description="A soldier hunts for hidden Ainu gold in Hokkaido" },
    @{ title="Chainsaw Man"; author="Tatsuki Fujimoto"; genres=@("Action","Horror","Shounen"); status="ongoing"; total_chapters=180; description="A boy merges with a devil to become a devil hunter" },
    @{ title="Spy x Family"; author="Tatsuya Endo"; genres=@("Action","Comedy","Shounen"); status="ongoing"; total_chapters=110; description="A spy builds a fake family for a secret mission" },
    @{ title="Tokyo Revengers"; author="Ken Wakui"; genres=@("Action","Drama","Shounen"); status="completed"; total_chapters=278; description="A man travels back in time to save his girlfriend" },
    @{ title="Kaiju No 8"; author="Naoya Matsumoto"; genres=@("Action","Sci-fi","Shounen"); status="ongoing"; total_chapters=110; description="A man gains the ability to transform into a kaiju" },
    @{ title="Blue Lock"; author="Muneyuki Kaneshiro"; genres=@("Sports","Shounen"); status="ongoing"; total_chapters=270; description="Strikers compete to become Japans best soccer player" },
    @{ title="Dandadan"; author="Yukinobu Tatsu"; genres=@("Action","Comedy","Supernatural"); status="ongoing"; total_chapters=180; description="A girl and boy encounter ghosts and aliens together" },
    @{ title="Frieren"; author="Kanehito Yamada"; genres=@("Adventure","Fantasy","Drama"); status="ongoing"; total_chapters=130; description="An elf mage reflects on life after defeating the demon king" },
    @{ title="Oshi no Ko"; author="Aka Akasaka"; genres=@("Drama","Mystery","Seinen"); status="ongoing"; total_chapters=160; description="A doctor reincarnates as the child of his idol" },
    @{ title="Sakamoto Days"; author="Yuto Suzuki"; genres=@("Action","Comedy","Shounen"); status="ongoing"; total_chapters=180; description="A legendary hitman retires to run a convenience store" },
    @{ title="Wind Breaker"; author="Satoru Nii"; genres=@("Action","Sports","Shounen"); status="ongoing"; total_chapters=110; description="A delinquent joins a crew that protects their town" }
)

Write-Host "========================================"
Write-Host "3. THEM $($mangas.Count) MANGA..."
Write-Host ""

$i = 1
foreach ($manga in $mangas) {
    Write-Host "[$i/$($mangas.Count)] Adding: $($manga.title)"
    try {
        $body = $manga | ConvertTo-Json
        $res = Invoke-RestMethod -Method POST "$BASE_URL/manga" `
            -ContentType "application/json" `
            -Headers $headers `
            -Body $body
        Write-Host "  OK: $($res.manga.id)"
    } catch {
        Write-Host "  LOI: $_"
    }
    $i++
}

Write-Host ""
Write-Host "========================================"
Write-Host "DANG XUAT..."
Invoke-RestMethod -Method POST "$BASE_URL/auth/logout" -Headers $headers
Write-Host "DA HOAN THANH $($mangas.Count) MANGA!"