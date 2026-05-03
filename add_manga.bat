@echo off
set TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImpvaG5AZ21haWwuY29tIiwiZXhwIjoxNzc3NzczNDA5LCJpYXQiOjE3Nzc3NzI4MDksInVzZXJfaWQiOiJ1c3JfMTc3NzMwMjE2NzQ0NzI5MzIwMCIsInVzZXJuYW1lIjoiam9obmRvZSJ9._nWgBK63esmFoWnpWSXyqsjxx2F-bFlnJbDhA8UoRnQ
echo ADDING 100 MANGA TO DATABASE
echo ========================================

REM ===== SHOUNEN =====
echo [1] Naruto
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Naruto\", \"author\": \"Kishimoto Masashi\", \"genres\": \"Action,Adventure,Shounen\", \"status\": \"completed\", \"total_chapters\": 700, \"description\": \"A young ninja seeks recognition and dreams of becoming Hokage\"}"

echo [2] Dragon Ball
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Dragon Ball\", \"author\": \"Toriyama Akira\", \"genres\": \"Action,Adventure,Shounen\", \"status\": \"completed\", \"total_chapters\": 519, \"description\": \"Son Goku's adventures searching for Dragon Balls\"}"

echo [3] Bleach
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Bleach\", \"author\": \"Kubo Tite\", \"genres\": \"Action,Supernatural,Shounen\", \"status\": \"completed\", \"total_chapters\": 686, \"description\": \"A teenager gains the powers of a Soul Reaper\"}"

echo [4] Attack on Titan
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Attack on Titan\", \"author\": \"Isayama Hajime\", \"genres\": \"Action,Drama,Fantasy,Shounen\", \"status\": \"completed\", \"total_chapters\": 139, \"description\": \"Humanity fights for survival against giant humanoid Titans\"}"

echo [5] My Hero Academia
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"My Hero Academia\", \"author\": \"Horikoshi Kohei\", \"genres\": \"Action,Superhero,Shounen\", \"status\": \"completed\", \"total_chapters\": 430, \"description\": \"A boy without powers in a world of superheroes\"}"

echo [6] Demon Slayer
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Demon Slayer\", \"author\": \"Gotouge Koyoharu\", \"genres\": \"Action,Supernatural,Shounen\", \"status\": \"completed\", \"total_chapters\": 205, \"description\": \"A boy becomes a demon slayer to save his sister\"}"

echo [7] Jujutsu Kaisen
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Jujutsu Kaisen\", \"author\": \"Akutami Gege\", \"genres\": \"Action,Supernatural,Shounen\", \"status\": \"ongoing\", \"total_chapters\": 260, \"description\": \"A boy swallows a cursed object and enters the world of jujutsu sorcerers\"}"

echo [8] Hunter x Hunter
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Hunter x Hunter\", \"author\": \"Togashi Yoshihiro\", \"genres\": \"Action,Adventure,Shounen\", \"status\": \"ongoing\", \"total_chapters\": 400, \"description\": \"A boy searches for his father by becoming a Hunter\"}"

echo [9] Fullmetal Alchemist
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Fullmetal Alchemist\", \"author\": \"Arakawa Hiromu\", \"genres\": \"Action,Adventure,Fantasy,Shounen\", \"status\": \"completed\", \"total_chapters\": 108, \"description\": \"Two brothers use alchemy to restore their bodies\"}"

echo [10] Fairy Tail
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Fairy Tail\", \"author\": \"Mashima Hiro\", \"genres\": \"Action,Adventure,Fantasy,Shounen\", \"status\": \"completed\", \"total_chapters\": 545, \"description\": \"Adventures of the Fairy Tail wizard guild\"}"

echo [11] Black Clover
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Black Clover\", \"author\": \"Tabata Yuki\", \"genres\": \"Action,Fantasy,Shounen\", \"status\": \"ongoing\", \"total_chapters\": 370, \"description\": \"A boy born without magic dreams of becoming the Wizard King\"}"

echo [12] Haikyuu
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Haikyuu!!\", \"author\": \"Furudate Haruichi\", \"genres\": \"Sports,Shounen\", \"status\": \"completed\", \"total_chapters\": 402, \"description\": \"A short boy pursues his dream of becoming a volleyball player\"}"

echo [13] Kuroko no Basket
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Kuroko no Basket\", \"author\": \"Fujimaki Tadatoshi\", \"genres\": \"Sports,Shounen\", \"status\": \"completed\", \"total_chapters\": 276, \"description\": \"A phantom sixth player joins a high school basketball team\"}"

echo [14] Toriko
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Toriko\", \"author\": \"Shimabukuro Mitsutoshi\", \"genres\": \"Action,Adventure,Shounen\", \"status\": \"completed\", \"total_chapters\": 396, \"description\": \"A gourmet hunter searches for the ultimate ingredients\"}"

echo [15] Sword Art Online
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Sword Art Online\", \"author\": \"Kawahara Reki\", \"genres\": \"Action,Fantasy,Shounen\", \"status\": \"ongoing\", \"total_chapters\": 30, \"description\": \"Players trapped in a virtual reality MMORPG\"}"

REM ===== SEINEN =====
echo [16] Berserk
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Berserk\", \"author\": \"Miura Kentaro\", \"genres\": \"Action,Dark Fantasy,Seinen\", \"status\": \"ongoing\", \"total_chapters\": 370, \"description\": \"A lone mercenary struggles in a dark medieval world\"}"

echo [17] Vagabond
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Vagabond\", \"author\": \"Inoue Takehiko\", \"genres\": \"Action,Historical,Seinen\", \"status\": \"ongoing\", \"total_chapters\": 327, \"description\": \"The life of legendary samurai Miyamoto Musashi\"}"

echo [18] Vinland Saga
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Vinland Saga\", \"author\": \"Yukimura Makoto\", \"genres\": \"Action,Historical,Seinen\", \"status\": \"ongoing\", \"total_chapters\": 200, \"description\": \"A Viking warrior seeks revenge and a land of peace\"}"

echo [19] Tokyo Ghoul
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Tokyo Ghoul\", \"author\": \"Ishida Sui\", \"genres\": \"Action,Horror,Seinen\", \"status\": \"completed\", \"total_chapters\": 179, \"description\": \"A college student becomes half-ghoul after a fatal encounter\"}"

echo [20] Gantz
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Gantz\", \"author\": \"Oku Hiroya\", \"genres\": \"Action,Sci-Fi,Seinen\", \"status\": \"completed\", \"total_chapters\": 383, \"description\": \"Dead people are forced to fight aliens in deadly missions\"}"

echo [21] Punpun
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Goodnight Punpun\", \"author\": \"Asano Inio\", \"genres\": \"Drama,Psychological,Seinen\", \"status\": \"completed\", \"total_chapters\": 147, \"description\": \"The dark coming-of-age story of a boy named Punpun\"}"

echo [22] Dungeon Meshi
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Dungeon Meshi\", \"author\": \"Ryoko Kui\", \"genres\": \"Adventure,Fantasy,Comedy,Seinen\", \"status\": \"completed\", \"total_chapters\": 97, \"description\": \"Adventurers cook and eat monsters in a dungeon\"}"

echo [23] Chainsaw Man
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Chainsaw Man\", \"author\": \"Fujimoto Tatsuki\", \"genres\": \"Action,Horror,Seinen\", \"status\": \"ongoing\", \"total_chapters\": 160, \"description\": \"A boy merges with his devil dog to become Chainsaw Man\"}"

echo [24] Fire Punch
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Fire Punch\", \"author\": \"Fujimoto Tatsuki\", \"genres\": \"Action,Drama,Seinen\", \"status\": \"completed\", \"total_chapters\": 83, \"description\": \"A boy with regeneration seeks revenge in a frozen world\"}"

echo [25] Biomega
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Biomega\", \"author\": \"Nihei Tsutomu\", \"genres\": \"Action,Sci-Fi,Seinen\", \"status\": \"completed\", \"total_chapters\": 48, \"description\": \"A synthetic human travels a post-apocalyptic world\"}"

echo [26] Monster
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Monster\", \"author\": \"Urasawa Naoki\", \"genres\": \"Thriller,Mystery,Seinen\", \"status\": \"completed\", \"total_chapters\": 162, \"description\": \"A doctor hunts down the patient he once saved\"}"

echo [27] 20th Century Boys
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"20th Century Boys\", \"author\": \"Urasawa Naoki\", \"genres\": \"Thriller,Mystery,Seinen\", \"status\": \"completed\", \"total_chapters\": 249, \"description\": \"Friends face a mysterious villain from their childhood\"}"

echo [28] Homunculus
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Homunculus\", \"author\": \"Yamamoto Hideo\", \"genres\": \"Drama,Psychological,Seinen\", \"status\": \"completed\", \"total_chapters\": 166, \"description\": \"A man gains the ability to see people's trauma\"}"

echo [29] Oyasumi Punpun
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"I Am a Hero\", \"author\": \"Hanazawa Kengo\", \"genres\": \"Horror,Seinen\", \"status\": \"completed\", \"total_chapters\": 296, \"description\": \"A manga artist survives a zombie apocalypse\"}"

echo [30] Solanin
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Solanin\", \"author\": \"Asano Inio\", \"genres\": \"Drama,Romance,Seinen\", \"status\": \"completed\", \"total_chapters\": 40, \"description\": \"Young adults navigate life and dreams in Tokyo\"}"

REM ===== SHOUJO =====
echo [31] Fruits Basket
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Fruits Basket\", \"author\": \"Takaya Natsuki\", \"genres\": \"Romance,Drama,Shoujo\", \"status\": \"completed\", \"total_chapters\": 136, \"description\": \"A girl lives with a family cursed to transform into zodiac animals\"}"

echo [32] Sailor Moon
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Sailor Moon\", \"author\": \"Takeuchi Naoko\", \"genres\": \"Romance,Magical Girl,Shoujo\", \"status\": \"completed\", \"total_chapters\": 52, \"description\": \"A schoolgirl transforms into a guardian to protect the world\"}"

echo [33] Cardcaptor Sakura
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Cardcaptor Sakura\", \"author\": \"CLAMP\", \"genres\": \"Magical Girl,Romance,Shoujo\", \"status\": \"completed\", \"total_chapters\": 50, \"description\": \"A girl must collect magical Clow Cards\"}"

echo [34] Ouran Host Club
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Ouran High School Host Club\", \"author\": \"Hatori Bisco\", \"genres\": \"Comedy,Romance,Shoujo\", \"status\": \"completed\", \"total_chapters\": 87, \"description\": \"A scholarship student accidentally joins an elite host club\"}"

echo [35] Skip Beat
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Skip Beat!\", \"author\": \"Nakamura Yoshiki\", \"genres\": \"Drama,Romance,Shoujo\", \"status\": \"ongoing\", \"total_chapters\": 300, \"description\": \"A girl enters showbiz to get revenge on her ex\"}"

echo [36] Kimi ni Todoke
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Kimi ni Todoke\", \"author\": \"Shiina Karuho\", \"genres\": \"Romance,Drama,Shoujo\", \"status\": \"completed\", \"total_chapters\": 124, \"description\": \"A misunderstood girl finds friendship and love\"}"

echo [37] Nana
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Nana\", \"author\": \"Yazawa Ai\", \"genres\": \"Drama,Romance,Shoujo\", \"status\": \"ongoing\", \"total_chapters\": 84, \"description\": \"Two girls named Nana meet and become roommates in Tokyo\"}"

echo [38] Clannad
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Clannad\", \"author\": \"Key\", \"genres\": \"Drama,Romance,Shoujo\", \"status\": \"completed\", \"total_chapters\": 105, \"description\": \"A delinquent meets a girl and discovers the meaning of family\"}"

echo [39] Itazura na Kiss
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Itazura na Kiss\", \"author\": \"Tada Kaoru\", \"genres\": \"Comedy,Romance,Shoujo\", \"status\": \"completed\", \"total_chapters\": 23, \"description\": \"A clumsy girl pursues the smartest boy in school\"}"

echo [40] Kaichou wa Maid-sama
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Kaichou wa Maid-sama!\", \"author\": \"Fujiwara Hiro\", \"genres\": \"Comedy,Romance,Shoujo\", \"status\": \"completed\", \"total_chapters\": 85, \"description\": \"A strict student council president secretly works at a maid cafe\"}"

REM ===== JOSEI =====
echo [41] Natsume Yuujinchou
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Natsume Yuujinchou\", \"author\": \"Midorikawa Yuki\", \"genres\": \"Supernatural,Slice of Life,Josei\", \"status\": \"ongoing\", \"total_chapters\": 115, \"description\": \"A boy can see spirits and inherits his grandmother's book of friends\"}"

echo [42] Chihayafuru
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Chihayafuru\", \"author\": \"Suetsugu Yuki\", \"genres\": \"Sports,Romance,Josei\", \"status\": \"completed\", \"total_chapters\": 246, \"description\": \"A girl pursues her dream of becoming the best karuta player\"}"

echo [43] Honey and Clover
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Honey and Clover\", \"author\": \"Umino Chica\", \"genres\": \"Drama,Romance,Josei\", \"status\": \"completed\", \"total_chapters\": 78, \"description\": \"Art students navigate love and life in college\"}"

echo [44] Paradise Kiss
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Paradise Kiss\", \"author\": \"Yazawa Ai\", \"genres\": \"Drama,Romance,Josei\", \"status\": \"completed\", \"total_chapters\": 42, \"description\": \"A high school girl becomes a model for fashion students\"}"

echo [45] Nodame Cantabile
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Nodame Cantabile\", \"author\": \"Ninomiya Tomoko\", \"genres\": \"Comedy,Romance,Josei\", \"status\": \"completed\", \"total_chapters\": 150, \"description\": \"Two music students with different personalities fall in love\"}"

REM ===== ISEKAI =====
echo [46] Re:Zero
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Re:Zero\", \"author\": \"Nagatsuki Tappei\", \"genres\": \"Fantasy,Isekai,Action\", \"status\": \"ongoing\", \"total_chapters\": 90, \"description\": \"A boy is transported to a fantasy world with the ability to respawn\"}"

echo [47] Overlord
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Overlord\", \"author\": \"Maruyama Kugane\", \"genres\": \"Fantasy,Isekai,Action\", \"status\": \"ongoing\", \"total_chapters\": 70, \"description\": \"A player is trapped in a game world as a powerful undead ruler\"}"

echo [48] That Time I Got Reincarnated as a Slime
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"That Time I Got Reincarnated as a Slime\", \"author\": \"Fuse\", \"genres\": \"Fantasy,Isekai,Comedy\", \"status\": \"ongoing\", \"total_chapters\": 110, \"description\": \"A man reincarnates as a slime in a fantasy world\"}"

echo [49] Konosuba
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"KonoSuba\", \"author\": \"Akatsuki Natsume\", \"genres\": \"Fantasy,Isekai,Comedy\", \"status\": \"ongoing\", \"total_chapters\": 70, \"description\": \"A boy and a useless goddess go on comedic adventures\"}"

echo [50] Mushoku Tensei
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Mushoku Tensei\", \"author\": \"Rifujin na Magonote\", \"genres\": \"Fantasy,Isekai,Adventure\", \"status\": \"ongoing\", \"total_chapters\": 90, \"description\": \"A jobless man reincarnates and lives his new life to the fullest\"}"

REM ===== ROMANCE/SLICE OF LIFE =====
echo [51] Toradora
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Toradora!\", \"author\": \"Takemiya Yuyuko\", \"genres\": \"Romance,Comedy,Slice of Life\", \"status\": \"completed\", \"total_chapters\": 50, \"description\": \"Two students help each other confess to their crushes\"}"

echo [52] Chuunibyou
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Chuunibyou demo Koi ga Shitai\", \"author\": \"Torako\", \"genres\": \"Romance,Comedy,Slice of Life\", \"status\": \"completed\", \"total_chapters\": 30, \"description\": \"A boy tries to leave his embarrassing past but meets a delusional girl\"}"

echo [53] Nisekoi
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Nisekoi\", \"author\": \"Komi Naoshi\", \"genres\": \"Romance,Comedy,Harem\", \"status\": \"completed\", \"total_chapters\": 229, \"description\": \"Two students from rival gangs must pretend to be a couple\"}"

echo [54] Kaguya-sama
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Kaguya-sama: Love is War\", \"author\": \"Akasaka Aka\", \"genres\": \"Romance,Comedy\", \"status\": \"completed\", \"total_chapters\": 281, \"description\": \"Two genius students try to make each other confess first\"}"

echo [55] Rent a Girlfriend
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Rent-a-Girlfriend\", \"author\": \"Miyajima Reiji\", \"genres\": \"Romance,Comedy\", \"status\": \"ongoing\", \"total_chapters\": 320, \"description\": \"A college student rents a girlfriend and gets into a complicated situation\"}"

echo [56] Horimiya
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Horimiya\", \"author\": \"HERO\", \"genres\": \"Romance,Slice of Life\", \"status\": \"completed\", \"total_chapters\": 122, \"description\": \"Two students discover their hidden sides and fall in love\"}"

echo [57] Ao Haru Ride
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Ao Haru Ride\", \"author\": \"Sakisaka Io\", \"genres\": \"Romance,Drama,Shoujo\", \"status\": \"completed\", \"total_chapters\": 49, \"description\": \"A girl reunites with her first love who has changed\"}"

echo [58] Your Lie in April
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Your Lie in April\", \"author\": \"Arakawa Naoshi\", \"genres\": \"Drama,Romance,Music\", \"status\": \"completed\", \"total_chapters\": 44, \"description\": \"A pianist meets a violinist who changes his world\"}"

echo [59] A Silent Voice
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"A Silent Voice\", \"author\": \"Oima Yoshitoki\", \"genres\": \"Drama,Romance,Slice of Life\", \"status\": \"completed\", \"total_chapters\": 62, \"description\": \"A boy seeks redemption after bullying a deaf classmate\"}"

echo [60] Takagi-san
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Karakai Jouzu no Takagi-san\", \"author\": \"Yamamoto Souichirou\", \"genres\": \"Romance,Comedy,Slice of Life\", \"status\": \"ongoing\", \"total_chapters\": 190, \"description\": \"A girl teases her classmate every day\"}"

REM ===== MYSTERY/THRILLER =====
echo [61] Death Note
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Death Note\", \"author\": \"Ohba Tsugumi\", \"genres\": \"Thriller,Mystery,Supernatural\", \"status\": \"completed\", \"total_chapters\": 108, \"description\": \"A student finds a notebook that can kill anyone\"}"

echo [62] The Promised Neverland
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"The Promised Neverland\", \"author\": \"Shirai Kaiu\", \"genres\": \"Thriller,Mystery,Horror\", \"status\": \"completed\", \"total_chapters\": 181, \"description\": \"Children in an orphanage discover a dark secret\"}"

echo [63] Erased
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Erased\", \"author\": \"Sanbe Kei\", \"genres\": \"Mystery,Thriller,Drama\", \"status\": \"completed\", \"total_chapters\": 44, \"description\": \"A man travels back in time to prevent a murder\"}"

echo [64] Billy Bat
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Billy Bat\", \"author\": \"Urasawa Naoki\", \"genres\": \"Mystery,Thriller,Historical\", \"status\": \"completed\", \"total_chapters\": 166, \"description\": \"A manga artist uncovers a conspiracy spanning centuries\"}"

echo [65] Q.E.D
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Q.E.D.\", \"author\": \"Motohashi Katou\", \"genres\": \"Mystery,Drama\", \"status\": \"ongoing\", \"total_chapters\": 100, \"description\": \"A genius student solves mysteries with his friend\"}"

REM ===== SCI-FI =====
echo [66] Akira
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Akira\", \"author\": \"Otomo Katsuhiro\", \"genres\": \"Sci-Fi,Action,Seinen\", \"status\": \"completed\", \"total_chapters\": 120, \"description\": \"A biker gang member gains powerful psychic abilities in Neo-Tokyo\"}"

echo [67] Planetes
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Planetes\", \"author\": \"Yukimura Makoto\", \"genres\": \"Sci-Fi,Drama,Seinen\", \"status\": \"completed\", \"total_chapters\": 26, \"description\": \"Space debris collectors in the near future\"}"

echo [68] Blame
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Blame!\", \"author\": \"Nihei Tsutomu\", \"genres\": \"Sci-Fi,Action,Seinen\", \"status\": \"completed\", \"total_chapters\": 66, \"description\": \"A man searches for a gene in a vast cyberpunk megastructure\"}"

echo [69] Dorohedoro
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Dorohedoro\", \"author\": \"Hayashida Q\", \"genres\": \"Action,Fantasy,Dark Comedy\", \"status\": \"completed\", \"total_chapters\": 167, \"description\": \"A man with a lizard head searches for the sorcerer who cursed him\"}"

echo [70] Biomeat Nectar
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Bio-Meat: Nectar\", \"author\": \"Fujisawa Yuki\", \"genres\": \"Horror,Sci-Fi\", \"status\": \"completed\", \"total_chapters\": 100, \"description\": \"Genetically engineered organisms escape and terrorize humans\"}"

REM ===== COMEDY =====
echo [71] Gintama
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Gintama\", \"author\": \"Sorachi Hideaki\", \"genres\": \"Comedy,Action,Sci-Fi,Shounen\", \"status\": \"completed\", \"total_chapters\": 704, \"description\": \"A samurai works odd jobs in an alien-occupied feudal Japan\"}"

echo [72] One Punch Man
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"One Punch Man\", \"author\": \"ONE\", \"genres\": \"Action,Comedy,Superhero\", \"status\": \"ongoing\", \"total_chapters\": 200, \"description\": \"A hero who can defeat any enemy with one punch seeks a worthy challenge\"}"

echo [73] Saiki K
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"The Disastrous Life of Saiki K.\", \"author\": \"Asou Shuuichi\", \"genres\": \"Comedy,Supernatural,Shounen\", \"status\": \"completed\", \"total_chapters\": 282, \"description\": \"A powerful psychic tries to live a normal life\"}"

echo [74] Mob Psycho 100
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Mob Psycho 100\", \"author\": \"ONE\", \"genres\": \"Action,Comedy,Supernatural\", \"status\": \"completed\", \"total_chapters\": 101, \"description\": \"A powerful psychic boy tries to live normally under his con-artist mentor\"}"

echo [75] Grand Blue
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Grand Blue Dreaming\", \"author\": \"Inoue Kenji\", \"genres\": \"Comedy,Slice of Life,Seinen\", \"status\": \"ongoing\", \"total_chapters\": 80, \"description\": \"A college student joins a diving club full of party-loving members\"}"

REM ===== HORROR =====
echo [76] Uzumaki
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Uzumaki\", \"author\": \"Ito Junji\", \"genres\": \"Horror,Mystery\", \"status\": \"completed\", \"total_chapters\": 19, \"description\": \"A town is cursed by spiral shapes\"}"

echo [77] Tomie
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Tomie\", \"author\": \"Ito Junji\", \"genres\": \"Horror,Mystery\", \"status\": \"completed\", \"total_chapters\": 20, \"description\": \"An immortal girl drives men to madness and murder\"}"

echo [78] Hellstar Remina
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Remina\", \"author\": \"Ito Junji\", \"genres\": \"Horror,Sci-Fi\", \"status\": \"completed\", \"total_chapters\": 6, \"description\": \"A planet devours other planets and heads toward Earth\"}"

echo [79] Another
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Another\", \"author\": \"Ayatsuji Yukito\", \"genres\": \"Horror,Mystery\", \"status\": \"completed\", \"total_chapters\": 26, \"description\": \"A transfer student encounters a deadly class curse\"}"

echo [80] Parasyte
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Parasyte\", \"author\": \"Iwaaki Hitoshi\", \"genres\": \"Horror,Sci-Fi,Seinen\", \"status\": \"completed\", \"total_chapters\": 64, \"description\": \"Alien parasites invade human bodies but one fails to reach the brain\"}"

REM ===== HISTORICAL =====
echo [81] Kingdom
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Kingdom\", \"author\": \"Hara Yasuhisa\", \"genres\": \"Action,Historical,Seinen\", \"status\": \"ongoing\", \"total_chapters\": 780, \"description\": \"A war orphan rises to become the greatest general in ancient China\"}"

echo [82] Lone Wolf and Cub
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Lone Wolf and Cub\", \"author\": \"Koike Kazuo\", \"genres\": \"Action,Historical,Seinen\", \"status\": \"completed\", \"total_chapters\": 142, \"description\": \"A ronin and his son travel a path of vengeance\"}"

echo [83] Rurouni Kenshin
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Rurouni Kenshin\", \"author\": \"Watsuki Nobuhiro\", \"genres\": \"Action,Historical,Shounen\", \"status\": \"completed\", \"total_chapters\": 255, \"description\": \"A legendary assassin seeks to atone for his past sins\"}"

echo [84] Slam Dunk
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Slam Dunk\", \"author\": \"Inoue Takehiko\", \"genres\": \"Sports,Comedy,Shounen\", \"status\": \"completed\", \"total_chapters\": 276, \"description\": \"A delinquent joins the basketball team to impress a girl\"}"

echo [85] Captain Tsubasa
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Captain Tsubasa\", \"author\": \"Takahashi Yoichi\", \"genres\": \"Sports,Shounen\", \"status\": \"ongoing\", \"total_chapters\": 130, \"description\": \"A boy dreams of winning the FIFA World Cup\"}"

REM ===== FANTASY =====
echo [86] Made in Abyss
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Made in Abyss\", \"author\": \"Tsukushi Akihito\", \"genres\": \"Adventure,Fantasy,Dark Fantasy\", \"status\": \"ongoing\", \"total_chapters\": 65, \"description\": \"A girl descends into a mysterious abyss to find her mother\"}"

echo [87] Magi
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Magi: The Labyrinth of Magic\", \"author\": \"Ohtaka Shinobu\", \"genres\": \"Action,Fantasy,Adventure,Shounen\", \"status\": \"completed\", \"total_chapters\": 369, \"description\": \"A boy travels with a djinn through dungeons and kingdoms\"}"

echo [88] The Rising of the Shield Hero
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"The Rising of the Shield Hero\", \"author\": \"Yusagi Aneko\", \"genres\": \"Fantasy,Isekai,Action\", \"status\": \"ongoing\", \"total_chapters\": 90, \"description\": \"A falsely accused hero rises to become the most powerful\"}"

echo [89] Claymore
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Claymore\", \"author\": \"Yagi Norihiro\", \"genres\": \"Action,Fantasy,Dark Fantasy\", \"status\": \"completed\", \"total_chapters\": 155, \"description\": \"Half-human half-monster warriors fight demons\"}"

echo [90] Record of Ragnarok
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Record of Ragnarok\", \"author\": \"Umemura Shinya\", \"genres\": \"Action,Fantasy,Mythology\", \"status\": \"ongoing\", \"total_chapters\": 80, \"description\": \"Humanity's greatest fighters battle against the gods\"}"

REM ===== MECHA =====
echo [91] Neon Genesis Evangelion
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Neon Genesis Evangelion\", \"author\": \"Sadamoto Yoshiyuki\", \"genres\": \"Mecha,Sci-Fi,Drama\", \"status\": \"completed\", \"total_chapters\": 97, \"description\": \"A boy pilots a giant mech to fight mysterious beings called Angels\"}"

echo [92] Code Geass
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Code Geass\", \"author\": \"Taniguchi Goro\", \"genres\": \"Mecha,Sci-Fi,Action\", \"status\": \"completed\", \"total_chapters\": 60, \"description\": \"An exiled prince gains the power to control anyone\"}"

echo [93] Gurren Lagann
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Tengen Toppa Gurren Lagann\", \"author\": \"Nakashima Kazuki\", \"genres\": \"Mecha,Action,Sci-Fi\", \"status\": \"completed\", \"total_chapters\": 27, \"description\": \"A boy drills through destiny in a giant mech\"}"

REM ===== MORE POPULAR =====
echo [94] Dandadan
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Dandadan\", \"author\": \"Tatsu Yukinobu\", \"genres\": \"Action,Supernatural,Comedy\", \"status\": \"ongoing\", \"total_chapters\": 150, \"description\": \"Two students investigate occult and alien encounters\"}"

echo [95] Blue Lock
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Blue Lock\", \"author\": \"Kaneshiro Muneyuki\", \"genres\": \"Sports,Shounen\", \"status\": \"ongoing\", \"total_chapters\": 270, \"description\": \"300 strikers compete to become Japan's best forward\"}"

echo [96] Spy x Family
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Spy x Family\", \"author\": \"Endo Tatsuya\", \"genres\": \"Comedy,Action,Slice of Life\", \"status\": \"ongoing\", \"total_chapters\": 100, \"description\": \"A spy creates a fake family for a mission\"}"

echo [97] Frieren
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Frieren: Beyond Journey's End\", \"author\": \"Yamada Kanehito\", \"genres\": \"Fantasy,Adventure,Slice of Life\", \"status\": \"ongoing\", \"total_chapters\": 120, \"description\": \"An elven mage reflects on her past companions after defeating the demon king\"}"

echo [98] Oshi no Ko
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Oshi no Ko\", \"author\": \"Akasaka Aka\", \"genres\": \"Drama,Mystery,Supernatural\", \"status\": \"ongoing\", \"total_chapters\": 150, \"description\": \"A doctor is reincarnated as the child of his idol\"}"

echo [99] Kaiju No.8
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Kaiju No. 8\", \"author\": \"Matsumoto Naoya\", \"genres\": \"Action,Sci-Fi,Shounen\", \"status\": \"ongoing\", \"total_chapters\": 110, \"description\": \"A man transforms into a kaiju while working as a monster cleaner\"}"

echo [100] Wind Breaker
curl -s -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Wind Breaker\", \"author\": \"Nii Satoru\", \"genres\": \"Action,Shounen\", \"status\": \"ongoing\", \"total_chapters\": 100, \"description\": \"A delinquent joins a gang that protects their town\"}"

echo.
echo ========================================
echo DONE! Added 100 manga to database
echo ========================================
pause