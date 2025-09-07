package utils

import (
    "fmt"
    "math/rand"
    "regexp"
    "strings"
    "time"

    "golang.org/x/text/transform"
    "golang.org/x/text/unicode/norm"
)

var (
    
    vnLast = []string{
        "Nguyen","Tran","Le","Pham","Hoang","Huynh","Phan","Vu","Vo","Dang","Bui","Do","Ho","Ngo","Duong","Ly","Mai","Truong","Doan","Cao","Chau","Ta","To","Ha","Kieu","Diep","Dinh","Dao","Luong","Lam","La","Vuong","Trieu","Quach","Quang","Phung","Ton","Kim","Han","Canh","Cung","Cai","Cam","Tieu","Trac","Trieu","Luu","Trinh","Phu","Phuoc","Gia","Ngoan","Viet","Thach","Diep","Quan","Quoc","Thai","Thao","Thanh","Thang","Tien","Phat","Hau","Khanh","Khang","Kien","Kiet","Long","Phuc","Phong","Son","Sang","Van","Vuong","Minh","Huy","Hieu","Hai","Hung","Tuan","Trung","An","Anh","Bac","Bao","Binh","Chi","Cong","Cuong","Dai","Dang2","Dao2","Diem","Dinh2","Doan2","Duc","Duong2","Dat","Giang","Ha2","Hang","Han2","Hanh","Hao","Hien","Hoa","Hoai","Hoang2","Hong","Huan","Huong","Huyen","Khoa","Le2","Linh","Loan","Long2","Luong2","Ly2","Mai2","Manh","Mi","My","Nam","Nga","Ngan","Ngoc","Nguyet","Nghia","Nhi","Nhan","Nhat","Nhu","Quyen","Quynh","Suong","Tam","Tan","Thach2","Thien","Thinh","Tho","Thu","Thuan","Thuc","Thuy","Tin","Tinh","Toan","Trang","Tri","Trinh2","Truc","Truong2","Tu","Tuyet","Tuyen","Vi","Vy","Xuan","Yen",
    }
 
    vnGiven1 = []string{
        "an","anh","bao","bang","binh","bich","cam","canh","chau","chi","cong","cuong","dai","dang","dao","diem","diep","dinh","doan","duc","dung","duong","dat","gia","giang","ha","hai","hang","han","hanh","hao","hau","hieu","hien","hoa","hoai","hoang","hong","huan","hung","huy","huong","huyen","khanh","khang","khoa","kien","kiet","kim","lam","lan","lang","le","linh","loan","long","luong","ly","mai","manh","mi","minh","my","nam","nga","ngan","ngat","ngoc","nguyet","nghia","nhi","nhan","nhat","nhu","phat","phuc","phuong","phuoc","phong","quan","quang","quoc","quyen","quynh","sang","son","suong","tam","tan","thach","thai","thanh","thao","thang","thien","thinh","tho","thu","thuan","thuc","thuy","tien","tin","tinh","toan","trang","trieu","tri","trinh","truc","truong","trung","tu","tuan","tuyet","tuyen","van","vi","viet","vy","xuan","yen",
    }
    vnGiven2 = []string{
        "an","anh","bao","binh","chau","chi","duy","dung","ha","hai","han","hanh","hieu","hien","hoa","hoai","hoang","hong","hung","huy","huong","huyen","khanh","khoa","kiet","kim","lan","linh","long","ly","mai","minh","my","nam","ngoc","nghia","nhat","nhu","phuc","phuong","phuoc","phong","quan","quang","quoc","quynh","sang","son","tam","tan","thai","thanh","thao","thang","thien","thinh","thu","thuy","tien","toan","trang","tri","trinh","truc","truong","trung","tu","tuan","van","viet","vy","xuan","yen",
    }
    vnMiddle = []string{"van","thi","thu","thanh","ngoc","huu","duc","thuy","minh","anh"}

    enFirst = []string{
        "Aaron","Abigail","Adam","Adrian","Aiden","Alan","Albert","Alex","Alexa","Alexander","Alice","Alicia","Amanda","Amber","Amy","Andrea","Andrew","Angela","Anna","Anthony","Arthur","Ava","Barbara","Benjamin","Betty","Beverly","Blake","Brandon","Brenda","Brian","Brittany","Brooke","Caleb","Cameron","Carl","Carla","Carlos","Carolyn","Carter","Catherine","Chad","Charles","Charlotte","Chloe","Christian","Christina","Christine","Christopher","Cindy","Claire","Cole","Colin","Connor","Courtney","Crystal","Daniel","David","Dawn","Deborah","Debra","Denise","Dennis","Diana","Diane","Donald","Donna","Dylan","Edward","Elijah","Elizabeth","Ella","Emily","Emma","Eric","Ethan","Eugene","Evan","Evelyn","Faith","Frank","Gabriel","Gary","George","Grace","Gregory","Hannah","Harold","Harry","Heather","Helen","Henry","Isabella","Jack","Jackson","Jacob","James","Jamie","Jane","Janet","Jasmine","Jason","Jean","Jeffrey","Jenna","Jennifer","Jeremy","Jerry","Jesse","Jessica","Jill","Joan","Joe","John","Johnny","Jonathan","Jordan","Jose","Joseph","Joshua","Joyce","Juan","Judy","Julia","Julie","Justin","Karen","Katherine","Kathleen","Kathryn","Kayla","Keith","Kelly","Kenneth","Kevin","Kim","Kimberly","Kyle","Laura","Lauren","Lawrence","Liam","Lillian","Linda","Lisa","Logan","Lori","Louis","Madison","Makayla","Marcus","Margaret","Maria","Marie","Marilyn","Mark","Martha","Mary","Mason","Matthew","Megan","Melissa","Michael","Michelle","Molly","Morgan","Nancy","Natalie","Nathan","Nicholas","Nicole","Noah","Norman","Oliver","Olivia","Pamela","Patricia","Patrick","Paul","Peter","Philip","Rachel","Ralph","Randy","Raymond","Rebecca","Richard","Robert","Roger","Ronald","Rose","Roy","Russell","Ryan","Samantha","Samuel","Sandra","Sara","Sarah","Scott","Sean","Sharon","Shawn","Shirley","Sophia","Stephanie","Stephen","Steven","Susan","Tammy","Taylor","Teresa","Terry","Theresa","Thomas","Tiffany","Timothy","Tyler","Victoria","Vincent","Virginia","Walter","Wayne","William","Willow","Zachary",
     
        "Aarav","Abel","Aisha","Alina","Amelia","Ariana","Asher","Aurora","Avery","Axel","Beatrice","Bianca","Brayden","Callum","Caroline","Cassidy","Cecilia","Colton","Cooper","Daisy","Delilah","Easton","Eleanor","Elena","Eli","Eliana","Elise","Emerson","Emmett","Everett","Eva","Felix","Finley","Fiona","Gianna","Grace","Grayson","Hadley","Hailey","Harper","Hazel","Hudson","Ian","Iris","Ivy","Jade","Jasper","Jaxson","Jonah","Josiah","Jude","Keegan","Kennedy","Kingston","Knox","Kylie","Laila","Leah","Leon","Levi","Lincoln","Lola","Luca","Lucas","Lucy","Luna","Maddox","Maya","Mila","Myles","Naomi","Nora","Nova","Nyla","Oakley","Paisley","Parker","Peyton","Quinn","Reagan","Riley","River","Roman","Rose","Ruby","Ryder","Sadie","Sage","Sawyer","Scarlett","Sebastian","Serena","Skylar","Sloane","Stella","Theodore","Tristan","Vera","Violet","Wesley","Weston","Willow","Xavier","Zara","Zoe",
    }
    enLast = []string{
        "Adams","Allen","Anderson","Bailey","Baker","Barnes","Bell","Bennett","Brooks","Brown","Bryant","Butler","Campbell","Carter","Clark","Collins","Cook","Cooper","Cox","Davis","Diaz","Edwards","Evans","Fisher","Flores","Foster","Garcia","Gomez","Gonzalez","Gray","Green","Griffin","Hall","Harris","Hayes","Henderson","Hernandez","Hill","Howard","Hughes","Jackson","James","Jenkins","Johnson","Jones","Kelly","King","Lee","Lewis","Long","Lopez","Martin","Martinez","Miller","Mitchell","Moore","Morgan","Morris","Murphy","Nelson","Parker","Patterson","Perez","Perry","Peterson","Phillips","Powell","Price","Ramirez","Reed","Reyes","Richardson","Rivera","Roberts","Robinson","Rodriguez","Rogers","Ross","Russell","Sanders","Scott","Simmons","Smith","Stewart","Taylor","Thomas","Thompson","Torres","Turner","Walker","Ward","Washington","Watson","White","Williams","Wilson","Wood","Wright","Young",
      
        "Abbott","Acevedo","Aguilar","Aguirre","Albert","Alexander","Alvarado","Alvarez","Andrews","Armstrong","Arnold","Austin","Avila","Banks","Barber","Barker","Bates","Beck","Becker","Bishop","Black","Blair","Boone","Bowers","Bowman","Boyd","Boyle","Bradley","Brady","Bremner","Briggs","Brock","Burke","Burns","Burton","Bush","Butcher","Cain","Caldwell","Carpenter","Carr","Carson","Casey","Castillo","Castro","Chambers","Chandler","Chapman","Christensen","Clarkson","Cline","Cobb","Coleman","Colon","Conner","Conrad","Cooke","Copeland","Crawford","Cross","Cruz","Cummings","Curry","Curtis","Dalton","Daniel","Daniels","Davidson","Dawson","Decker","Delgado","Dixon","Dominquez","Donovan","Doyle","Drake","Dunn","Eaton","Elliott","Ellis","Erickson","Estrada","Farrell","Faulkner","Ferguson","Fields","Figueroa","Fleming","Ford","Fowler","Fox","Francis","Franklin","Fuller","Gallagher","Galloway","Gardner","Garner","George","Gibson","Gilbert","Gill","Glover","Goodman","Goodwin","Gordon","Graham","Grant","Graves","Greer","Griffith","Gross","Guerra","Guerrero","Gutierrez","Guzman","Hale","Hanson","Hardy","Harper","Harrington","Hart","Harvey","Hawkins","Haynes","Hernadez","Hines","Hodges","Hoffman","Hogan","Holland","Holmes","Holt","Hopkins","Horton","Houston","Howard","Howell","Hubbard","Hudson","Humphrey","Hunt","Hunter","Ingram","Jared","Jefferson","Johns","Joseph","Kelley","Kemp","Kennedy","Klein","Knight","Lamb","Lambert","Lane","Larson","Lawson","Leonard","Lindsey","Little","Livingston","Lloyd","Logan","Lucas","Lynch","Mack","Mann","Manning","Marsh","Marshall","Mason","Matthews","May","McCarthy","McCormick","McCoy","McDaniel","McDonald","McGee","McGuire","McKenzie","McKinney","McLaughlin","Mendez","Mendoza","Miles","Mills","Miranda","Mitchells","Monroe","Montgomery","Moon","Morales","Moran","Moss","Mueller","Murray","Navarro","Neal","Newman","Nguyen","Nichols","Nielsen","Nunez","Obrien","Oliver","Ortega","Ortiz","Osborne","Page","Palmer","Park","Parks","Parsons","Payne","Pearson","Pena","Perkins","Peters","Phelps","Pierce","Pittman","Porter","Potter","Pratt","Preston","Price2","Quinn","Ramsey","Ramos","Randall","Ray","Reeves","Reid","Rios","Robbins","Rojas","Roman","Rose2","Ross2","Rowe","Ruiz","Salazar","Sanchez","Sandoval","Saunders","Schmidt","Schneider","Schultz","Schwartz","Sharp","Shelton","Sherman","Silva","Sims","Singleton","Slater","Sloan","Soto","Spears","Spencer","Stanley","Stokes","Stone","Sullivan","Summers","Sutton","Swanson","Sweeney","Tate","Todd","Tucker","Tyler","Valdez","Vargas","Vasquez","Vega","Velasquez","Villarreal","Wade","Wagner","Walsh","Walters","Weaver","Weber","Welch","Wells","Wheeler","Whitaker","Wilkins","Wilkinson","Williamson","Willis","Wong","Woods","Workman","Yates","Young2","Zimmerman",
    }
)

func init() { rand.Seed(time.Now().UnixNano()) }


func slugifyEmailLocal(s string) string {
    t := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
        return r >= 0x0300 && r <= 0x036F 
    }), norm.NFC)
    out, _, _ := transform.String(t, s)
    out = strings.ToLower(out)
  
    out = strings.ReplaceAll(out, " ", "")
    out = strings.ReplaceAll(out, ".", "")

    re := regexp.MustCompile(`[^a-z0-9]`)
    out = re.ReplaceAllString(out, "")
    out = strings.TrimSpace(out)
    if out == "" {
        out = "user"
    }
  
    return out
}

func randomChoice(a []string) string { return a[rand.Intn(len(a))] }

func GenerateRandomLocalPart() string {
   
    var full string
    if rand.Intn(2) == 0 {
      
        last := randomChoice(vnLast)
        var parts []string
        if rand.Intn(2) == 0 { 
            parts = []string{randomChoice(vnMiddle)}
        }
      
        given := randomChoice(vnGiven1)
        if rand.Intn(2) == 0 { 
            given = given + " " + randomChoice(vnGiven2)
        }
        full = last + " " + strings.Join(append(parts, given), " ")
    } else {
       
        full = randomChoice(enFirst) + " " + randomChoice(enLast)
    }
    base := slugifyEmailLocal(full)
    suffix := 1 + rand.Intn(999)
    return fmt.Sprintf("%s%d", base, suffix)
}

func GenerateRandomEmail() (string, bool) {
    ds := GetDomains()
    if len(ds) == 0 {
        return "", false
    }
    
    for i := 0; i < 20; i++ {
        local := GenerateRandomLocalPart()
        domain := ds[rand.Intn(len(ds))]
        email := local + "@" + domain
        if !IsBlacklisted(email) {
            AddToBlacklist(email)
            return email, true
        }
    }
    return "", false
}
