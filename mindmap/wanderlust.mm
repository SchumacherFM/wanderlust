<map version="1.0.1">
<!-- To view this file, download free mind mapping software FreeMind from http://freemind.sourceforge.net -->
<node BACKGROUND_COLOR="#999900" CREATED="1408226058402" ID="ID_513394138" MODIFIED="1408231819235">
<richcontent TYPE="NODE"><html>
  <head>
    
  </head>
  <body>
    <p style="text-align: center">
      <b>Wanderlust</b>
    </p>
    <p style="text-align: center">
      Cache Warmer
    </p>
    <p style="text-align: center">
      with priorities
    </p>
  </body>
</html></richcontent>
<node CREATED="1408226077809" ID="ID_1082254437" MODIFIED="1410908421899" POSITION="right">
<richcontent TYPE="NODE"><html>
  <head>
    
  </head>
  <body>
    <p>
      <b>Picnic</b>
    </p>
    <p>
      The web interface (Angular.js)
    </p>
    <p>
      - Watches DB for new URLs and already processed.
    </p>
    <p>
      - Pushes results from DB to the user via Websockets
    </p>
    <p>
      - Handles configurations
    </p>
  </body>
</html></richcontent>
<cloud COLOR="#ff9966"/>
<node CREATED="1408228814250" ID="ID_1818745543" MODIFIED="1408665089248" TEXT="Actions">
<cloud COLOR="#ffcc66"/>
<node CREATED="1408228827778" ID="ID_79040744" MODIFIED="1408228835639" TEXT="Start the Wanderer"/>
<node CREATED="1408228842031" ID="ID_1123195613" MODIFIED="1408229170111" TEXT="Start the Brotzeit"/>
<node CREATED="1408229470319" ID="ID_1737173709" MODIFIED="1408229482405" TEXT="Purge URLs collected from Brotzeit by provisioner"/>
<node CREATED="1408229490288" ID="ID_1374484413" MODIFIED="1408229493396" TEXT="Purge all data"/>
<node CREATED="1408665089999" ID="ID_686051620" MODIFIED="1408665139322" TEXT="Configure internal Cron Scheduler for each or all provisioners"/>
<node CREATED="1410907631357" ID="ID_1779739970" MODIFIED="1411815595866" TEXT="User Management">
<node CREATED="1411815596742" ID="ID_1529916533" MODIFIED="1411815598752" TEXT="Login"/>
<node CREATED="1411815599597" ID="ID_384392676" MODIFIED="1411815601504" TEXT="Logout"/>
<node CREATED="1411815603142" ID="ID_92525725" MODIFIED="1411815607944" TEXT="Change Password"/>
<node CREATED="1411815609229" ID="ID_1362166790" MODIFIED="1411815614072" TEXT="Create new User"/>
</node>
</node>
<node CREATED="1408229245604" ID="ID_1214455075" MODIFIED="1408229284099" TEXT="Views">
<cloud COLOR="#ffff66"/>
<node CREATED="1408227170612" ID="ID_431038421" MODIFIED="1408227176107" TEXT="View performance metrics"/>
<node CREATED="1408229209703" ID="ID_919530530" MODIFIED="1408229228214" TEXT="View all URLs from the provisioners"/>
<node CREATED="1408227074344" ID="ID_290075585" MODIFIED="1408227090765" TEXT="View Amount of Running Background Workers"/>
<node CREATED="1408666866219" ID="ID_1483591538" MODIFIED="1408666911942" TEXT="View of 404/500, etc errors"/>
<node CREATED="1408666887730" ID="ID_103649480" MODIFIED="1408666898430" TEXT="View of other errors"/>
</node>
<node CREATED="1408229253107" ID="ID_636372056" MODIFIED="1408229301345" TEXT="Configurations">
<cloud COLOR="#cccc00"/>
<node CREATED="1408227096764" ID="ID_1917687630" MODIFIED="1408229329570" TEXT="Provisioners and their access data"/>
<node CREATED="1408229333431" ID="ID_920425934" MODIFIED="1408229394680" TEXT="Wanderer: Concurrency level"/>
<node CREATED="1408229358645" ID="ID_454668395" MODIFIED="1408229419886" TEXT="Brotzeit: Concurrency level"/>
<node CREATED="1411815422138" ID="ID_931734721" MODIFIED="1411815428181" TEXT="Authentication">
<node CREATED="1411815430195" ID="ID_1872763293" MODIFIED="1411815493011" TEXT="AUTH_LEVEL_IGNORE we don&apos;t need the user in this handler"/>
<node CREATED="1411815449939" ID="ID_228901137" MODIFIED="1411815502283" TEXT="AUTH_LEVEL_CHECK prefetch user, doesn&apos;t matter if not logged in"/>
<node CREATED="1411815455057" ID="ID_1156768440" MODIFIED="1411815545504" TEXT="AUTH_LEVEL_LOGIN_WAIT user required, 412 if not precondition failed (log in missing)"/>
<node CREATED="1411815458418" ID="ID_624971482" MODIFIED="1411815553439" TEXT="AUTH_LEVEL_LOGIN user required, 401 if not available"/>
<node CREATED="1411815461618" ID="ID_1557802446" MODIFIED="1411815564724" TEXT="AUTH_LEVEL_ADMIN admin required, 401 if no user, 403 if not admin"/>
</node>
</node>
<node CREATED="1411815198758" ID="ID_482878418" MODIFIED="1411815314750" TEXT="Routes">
<cloud COLOR="#cc9900"/>
<node CREATED="1411815034801" ID="ID_551665319" MODIFIED="1411815038496" TEXT="/ GET">
<node CREATED="1411815038497" ID="ID_776506454" MODIFIED="1411815085489" TEXT="Redirect dashboard/"/>
</node>
<node CREATED="1411815069576" ID="ID_1267062929" MODIFIED="1411815074639" TEXT="/dashboard">
<node CREATED="1411815074640" ID="ID_1308184647" MODIFIED="1411815081146" TEXT="Redirect to dashboard/"/>
</node>
<node CREATED="1411815095447" ID="ID_354439861" MODIFIED="1411815096472" TEXT="/favicon.ico"/>
<node CREATED="1411815105621" ID="ID_370226543" MODIFIED="1411815107197" TEXT="/">
<node CREATED="1411815107199" ID="ID_1089694666" MODIFIED="1411815124762" TEXT="Gzrice Box, JS, CSS, HTML, Images"/>
</node>
<node CREATED="1408431999379" ID="ID_1149314733" MODIFIED="1408432599270" TEXT="REST API">
<cloud COLOR="#cccc00"/>
<node CREATED="1411814820587" ID="ID_1236673550" MODIFIED="1411814822732" TEXT="auth">
<node CREATED="1411814865912" ID="ID_1576528255" MODIFIED="1411814897866" TEXT="/ GET session info"/>
<node CREATED="1411814871201" ID="ID_1190409325" MODIFIED="1411814910218" TEXT="/ POST Login"/>
<node CREATED="1411814875345" ID="ID_1955078901" MODIFIED="1411814914113" TEXT="/ DELETE Logout"/>
</node>
<node CREATED="1411814833659" ID="ID_236936447" MODIFIED="1411814835080" TEXT="users">
<node CREATED="1411814928894" ID="ID_1082572905" MODIFIED="1411814933705" TEXT="/ GET collection"/>
</node>
<node CREATED="1411814847826" ID="ID_1595602320" MODIFIED="1411814849036" TEXT="sysinfo">
<node CREATED="1411814943221" ID="ID_690204432" MODIFIED="1411814969662" TEXT="/ GET Goroutines, Wanderers, Brotzeit, Prov."/>
</node>
<node CREATED="1411814978995" ID="ID_1719569356" MODIFIED="1411814980166" TEXT="brotzeit">
<node CREATED="1411814994386" ID="ID_722052997" MODIFIED="1411814997301" TEXT="@todo"/>
</node>
<node CREATED="1411815006698" ID="ID_1527342828" MODIFIED="1411815007826" TEXT="wanderer">
<node CREATED="1411815007827" ID="ID_1963620825" MODIFIED="1411815010651" TEXT="@todo"/>
</node>
<node CREATED="1411815023209" ID="ID_283685973" MODIFIED="1411815024465" TEXT="provisioners">
<node CREATED="1411815024467" ID="ID_1076841306" MODIFIED="1411815026372" TEXT="@todo"/>
</node>
</node>
<node CREATED="1411815315649" ID="ID_199482069" MODIFIED="1411815319808" TEXT="Middleware">
<node CREATED="1411815319809" ID="ID_439476076" MODIFIED="1411815323704" TEXT="CORS"/>
<node CREATED="1411815329128" ID="ID_764157888" MODIFIED="1411815339897" TEXT="GZIP auto detection"/>
</node>
</node>
</node>
<node CREATED="1408226087413" ID="ID_1865183383" MODIFIED="1408665380632" POSITION="left" STYLE="bubble">
<richcontent TYPE="NODE"><html>
  <head>
    
  </head>
  <body>
    <p>
      <b>The Wanderer</b>
    </p>
    <p>
      - Reads new URLs from the DB, fetches them and saves the results into the DB.
    </p>
    <p>
      - Marks fetched URLs as already processed
    </p>
    <p>
      - Nearly unlimited goroutines
    </p>
    <p>
      - r/w and watch to DB
    </p>
  </body>
</html></richcontent>
<cloud COLOR="#0099ff"/>
<node CREATED="1408226766578" ID="ID_1513397460" MODIFIED="1408226865717" TEXT="Behaviour">
<node CREATED="1408226805365" ID="ID_1125581192" MODIFIED="1408226818214" TEXT="Gets list of URLs and dispatches them to the workers"/>
<node CREATED="1408226826134" ID="ID_1130130075" MODIFIED="1408226844094" TEXT="Returns Time for downloading and size for each URL"/>
</node>
<node CREATED="1408226783228" ID="ID_20049385" MODIFIED="1408232499822" TEXT="Worker Pool">
<icon BUILTIN="group"/>
<node CREATED="1408226111974" ID="ID_232571995" MODIFIED="1408226719515" TEXT="HTTP Worker"/>
<node CREATED="1408226151864" ID="ID_597746027" MODIFIED="1408226719515" TEXT="HTTP Worker"/>
<node CREATED="1408226158615" ID="ID_1945888817" MODIFIED="1408226719516" TEXT="HTTP Worker"/>
</node>
</node>
<node CREATED="1408226129104" ID="ID_815610302" MODIFIED="1408665047859" POSITION="left">
<richcontent TYPE="NODE"><html>
  <head>
    
  </head>
  <body>
    <p>
      <b>Provisioner</b>
    </p>
    <p>
      - List of modules which provides the URLs
    </p>
    <p>
      - Maintain Blacklist of URLs (e.g. no checkout)
    </p>
    <p>
      - Configure Priorities
    </p>
    <p>
      - Multiple Accounts of GA, Piwik, REST endpoints, etc
    </p>
  </body>
</html></richcontent>
<cloud COLOR="#99ff99"/>
<node CREATED="1408226282839" ID="ID_1071132970" MODIFIED="1408226471322" TEXT="Prio">
<icon BUILTIN="full-1"/>
<node CREATED="1408226173643" ID="ID_1103780609" MODIFIED="1408226897113" TEXT="Google Analytics">
<cloud COLOR="#33ff33"/>
<node CREATED="1408226189799" ID="ID_1637906665" MODIFIED="1408226225164" TEXT="Top 100 URLs last 2 weeks"/>
<node CREATED="1408226204381" ID="ID_1439236001" MODIFIED="1408226230063" TEXT="Top 200 sold products last 2 weeks"/>
</node>
<node CREATED="1408226234686" ID="ID_718481096" MODIFIED="1408226959610" TEXT="Piwik">
<cloud COLOR="#00ff99"/>
<node CREATED="1408226189799" ID="ID_1581195060" MODIFIED="1408229720906" TEXT="Top 100 URLs last x weeks"/>
<node CREATED="1408226204381" ID="ID_1100029901" MODIFIED="1408229724505" TEXT="Top 200 sold products last x weeks"/>
</node>
</node>
<node CREATED="1408226268878" ID="ID_1462097628" MODIFIED="1408226475394" TEXT="Prio">
<icon BUILTIN="full-2"/>
<node CREATED="1408226327769" ID="ID_120120562" MODIFIED="1408229524001">
<richcontent TYPE="NODE"><html>
  <head>
    
  </head>
  <body>
    <p>
      <b>JSON REST Endpoints</b>
    </p>
    <p>
      with/without OAuth config
    </p>
  </body>
</html></richcontent>
<cloud COLOR="#66ffcc"/>
<node CREATED="1408226343226" ID="ID_1516046397" MODIFIED="1408226351662" TEXT="e.g. Top 200 Sold products in Magento"/>
<node CREATED="1408226358864" ID="ID_139686099" MODIFIED="1408226375517" TEXT="e.g. CMS Pages"/>
<node CREATED="1408229744204" ID="ID_562106539" MODIFIED="1408229771084" TEXT="all other non-important products/pages"/>
<node CREATED="1408229832629" ID="ID_296927528" MODIFIED="1408229842489" TEXT="products/pages with the last cleared cache"/>
<node CREATED="1408232183237" ID="ID_1344362494" MODIFIED="1408232304006">
<richcontent TYPE="NODE"><html>
  <head>
    
  </head>
  <body style="text-align: left">
    <p>
      <b>JSON Format</b>
    </p>
    <pre style="font-size: 10px">{
    &quot;items&quot;: [
        {
            &quot;url&quot;: &quot;http://www.a.com/catalog/product/view/123&quot;
        },
        {
            &quot;url&quot;: &quot;http://www.a.com/about-us&quot;
        }
    ]
}</pre>
  </body>
</html></richcontent>
</node>
<node CREATED="1408321951235" ID="ID_1675573368" MODIFIED="1408321992192" TEXT="The order of the URLs in the JSON feed decides the loading order"/>
</node>
</node>
<node CREATED="1408226392635" ID="ID_1582515810" MODIFIED="1408226479578" TEXT="Prio">
<icon BUILTIN="full-3"/>
<node CREATED="1408226401120" ID="ID_1298495273" MODIFIED="1408226403356" TEXT="sitemap.xml"/>
</node>
<node CREATED="1408226615296" ID="ID_1295963359" MODIFIED="1408226622560" TEXT="Prio">
<icon BUILTIN="full-4"/>
<node CREATED="1408226628675" ID="ID_656036312" MODIFIED="1408226637556" TEXT="Manually added URLs"/>
</node>
<node CREATED="1408227537620" ID="ID_993462071" MODIFIED="1408227553916" TEXT="Prio">
<icon BUILTIN="full-0"/>
<node CREATED="1408227557645" ID="ID_1736240943" MODIFIED="1408229633490">
<richcontent TYPE="NODE"><html>
  <head>
    
  </head>
  <body>
    <p>
      Crawl all the remaining URLs of a website
    </p>
    <p>
      If configured
    </p>
  </body>
</html></richcontent>
</node>
</node>
</node>
<node CREATED="1408227030201" ID="ID_310886734" MODIFIED="1410907338502" POSITION="right">
<richcontent TYPE="NODE"><html>
  <head>
    
  </head>
  <body>
    <p>
      <b>The Rucksack</b>
    </p>
    <p>
      Database for storage. Compiled into.
    </p>
  </body>
</html></richcontent>
<cloud COLOR="#ccccff"/>
<node CREATED="1408321536418" ID="ID_1460198452" MODIFIED="1408321667611" TEXT="https://code.google.com/p/leveldb-go/"/>
<node CREATED="1408321580596" ID="ID_1138239741" MODIFIED="1408321924122" TEXT="https://github.com/cznic/kv"/>
<node CREATED="1408321633191" ID="ID_1710461451" MODIFIED="1408321634751" TEXT="https://github.com/boltdb/bolt"/>
<node CREATED="1408321676751" ID="ID_1351308280" MODIFIED="1408321677811" TEXT="https://github.com/syndtr/goleveldb"/>
<node CREATED="1408321894169" ID="ID_681319205" MODIFIED="1410907371304" TEXT="https://github.com/HouzuoGuo/tiedot">
<node CREATED="1410907415803" ID="ID_146269911" MODIFIED="1410907428823" TEXT="Stores JSON, like Mongodb ;-)"/>
</node>
</node>
<node CREATED="1408227268529" ID="ID_1662394047" MODIFIED="1408232458793" POSITION="left">
<richcontent TYPE="NODE"><html>
  <head>
    
  </head>
  <body>
    <p>
      <b>Brotzeit</b>
    </p>
    <p>
      -&#160;Fetches all the URLs from the Provisioner and stores them in the DB
    </p>
    <p>
      - Automatically fills the DB as soon as access data is provided to a provisioner or Enable: yes/no
    </p>
    <p>
      - Goroutine
    </p>
  </body>
</html></richcontent>
<cloud COLOR="#cc00cc"/>
</node>
<node CREATED="1408231473795" ID="ID_1639913958" MODIFIED="1408231477785" POSITION="right" TEXT="General Stuff">
<node CREATED="1408231437313" ID="ID_1943784766" MODIFIED="1408231442250" TEXT="Future Features">
<node CREATED="1408231445390" ID="ID_1320540073" MODIFIED="1408231459146" TEXT="Distributed Wanderers"/>
<node CREATED="1408332810790" ID="ID_1020350825" MODIFIED="1408332840229" TEXT="REST API in Picnic for the Rucksack"/>
<node CREATED="1408332893227" ID="ID_79108366" MODIFIED="1408332958913" TEXT="Extend Provisioner JSON to add NumOfRequests and Concurrency level"/>
</node>
<node CREATED="1408231265873" ID="ID_525514306" MODIFIED="1408231279512" TEXT="The only write access to the HDD is for saving the DB"/>
<node CREATED="1408231312099" ID="ID_1235551636" MODIFIED="1408231334046" TEXT="Wanderlust can run anywhere"/>
<node CREATED="1408231240797" ID="ID_225989754" MODIFIED="1410908446002" TEXT="All CSS, JS, HTML files are compiled into (via gz.rice)"/>
<node CREATED="1408321328218" ID="ID_970423774" MODIFIED="1411815950783">
<richcontent TYPE="NODE"><html>
  <head>
    
  </head>
  <body>
    CLI Settings
  </body>
</html>
</richcontent>
<node CREATED="1411815801223" ID="ID_285526862" MODIFIED="1411815844379" TEXT="Name:  &quot;picnic-listen-address,pla&quot;,&#xa;Value: &quot;127.0.0.1:3008&quot;,&#xa;Usage: &quot;IP:Port address for picnic dashboard to listen on&quot;,&#xa;"/>
<node CREATED="1411815860132" ID="ID_870649776" MODIFIED="1411815890159">
<richcontent TYPE="NODE"><html>
  <head>
    
  </head>
  <body>
    <p>
      Name:&#160;&#160;&quot;picnic-pem-dir,ppd&quot;,
    </p>
    <p>
      Value: &quot;&quot;,
    </p>
    <p>
      Usage: &quot;Directory to store the .pem certificates. If empty will throw it somewhere in the system. If provided file names must be cert.pem and key.pem!&quot;,
    </p>
  </body>
</html>
</richcontent>
</node>
<node CREATED="1411815951495" ID="ID_1081447868" MODIFIED="1411815983675" TEXT="Name:  &quot;rucksack-dir,rd&quot;,&#xa;Value: &quot;&quot;, &#xa;Usage: &quot;Storage directory of the rucksack files. If empty then /tmp/random directory will be used.&quot;, "/>
<node CREATED="1411816010949" ID="ID_1952015664" MODIFIED="1411816042016">
<richcontent TYPE="NODE"><html>
  <head>
    
  </head>
  <body>
    <p>
      Name:&#160;&#160;&quot;logFile,lf&quot;,
    </p>
    <p>
      Value: &quot;&quot;,
    </p>
    <p>
      Usage: &quot;Log to file or if empty to os.Stderr&quot;,
    </p>
  </body>
</html>
</richcontent>
</node>
<node CREATED="1411816050513" ID="ID_147859593" MODIFIED="1411816065934" TEXT="Name:  &quot;logLevel,ll&quot;,&#xa;Value: &quot;&quot;,&#xa;Usage: &quot;Log level: debug, info, notice, warning, error, critical, alert, emergency. Default: debug&quot;,&#xa;"/>
</node>
<node CREATED="1410907657982" ID="ID_1984817677" MODIFIED="1410908393960" TEXT="Maybe: Wanderer, Brotzeit and Rucksack communicate via Pub/Sub (bullhorn lib)"/>
<node CREATED="1410908458912" ID="ID_421948385" MODIFIED="1410908491476" TEXT="Only accessible via SSL"/>
</node>
</node>
</map>
