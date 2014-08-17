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
</html>
</richcontent>
<node CREATED="1408226077809" ID="ID_1082254437" MODIFIED="1408231745357" POSITION="right">
<richcontent TYPE="NODE"><html>
  <head>
    
  </head>
  <body>
    <p>
      <b>Picnic</b>
    </p>
    <p>
      The web interface
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
</html>
</richcontent>
<cloud COLOR="#ff9966"/>
<node CREATED="1408228814250" ID="ID_1818745543" MODIFIED="1408228910352" TEXT="Actions">
<cloud COLOR="#ffcc66"/>
<node CREATED="1408228827778" ID="ID_79040744" MODIFIED="1408228835639" TEXT="Start the Wanderer"/>
<node CREATED="1408228842031" ID="ID_1123195613" MODIFIED="1408229170111" TEXT="Start the Brotzeit"/>
<node CREATED="1408229470319" ID="ID_1737173709" MODIFIED="1408229482405" TEXT="Purge URLs collected from Brotzeit by provisioner"/>
<node CREATED="1408229490288" ID="ID_1374484413" MODIFIED="1408229493396" TEXT="Purge all data"/>
</node>
<node CREATED="1408229245604" ID="ID_1214455075" MODIFIED="1408229284099" TEXT="Views">
<cloud COLOR="#ffff66"/>
<node CREATED="1408227170612" ID="ID_431038421" MODIFIED="1408227176107" TEXT="View performance metrics"/>
<node CREATED="1408229209703" ID="ID_919530530" MODIFIED="1408229228214" TEXT="View all URLs from the provisioners"/>
<node CREATED="1408227074344" ID="ID_290075585" MODIFIED="1408227090765" TEXT="View Amount of Running Background Workers"/>
</node>
<node CREATED="1408229253107" ID="ID_636372056" MODIFIED="1408229301345" TEXT="Configurations">
<cloud COLOR="#cccc00"/>
<node CREATED="1408227096764" ID="ID_1917687630" MODIFIED="1408229329570" TEXT="Provisioners and their access data"/>
<node CREATED="1408229333431" ID="ID_920425934" MODIFIED="1408229394680" TEXT="Wanderer: Concurrency level"/>
<node CREATED="1408229358645" ID="ID_454668395" MODIFIED="1408229419886" TEXT="Brotzeit: Concurrency level"/>
</node>
</node>
<node CREATED="1408226087413" ID="ID_1865183383" MODIFIED="1408228749873" POSITION="left" STYLE="bubble">
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
      - Goroutine, r/w and watch to DB
    </p>
  </body>
</html>
</richcontent>
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
<node CREATED="1408226129104" ID="ID_815610302" MODIFIED="1408231856577" POSITION="left">
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
  </body>
</html>
</richcontent>
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
</html>
</richcontent>
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
</html>
</richcontent>
</node>
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
</html>
</richcontent>
</node>
</node>
</node>
<node CREATED="1408227030201" ID="ID_310886734" MODIFIED="1408228539032" POSITION="right">
<richcontent TYPE="NODE"><html>
  <head>
    
  </head>
  <body>
    <p>
      <b>The Rucksack</b>
    </p>
    <p>
      KV? Database for storage. No 3rd party. Compiled into.
    </p>
  </body>
</html>
</richcontent>
<cloud COLOR="#ccccff"/>
<node CREATED="1408231284740" ID="ID_25161220" MODIFIED="1408231295352" TEXT="Write access to HDD only for saving"/>
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
</html>
</richcontent>
<cloud COLOR="#cc00cc"/>
</node>
<node CREATED="1408231473795" ID="ID_1639913958" MODIFIED="1408231477785" POSITION="right" TEXT="General Stuff">
<node CREATED="1408231437313" ID="ID_1943784766" MODIFIED="1408231442250" TEXT="Future Features">
<node CREATED="1408231445390" ID="ID_1320540073" MODIFIED="1408231459146" TEXT="Distributed Wanderers"/>
</node>
<node CREATED="1408231265873" ID="ID_525514306" MODIFIED="1408231279512" TEXT="The only write access to the HDD is for saving the DB"/>
<node CREATED="1408231312099" ID="ID_1235551636" MODIFIED="1408231334046" TEXT="Wanderlust can run anywhere"/>
<node CREATED="1408231240797" ID="ID_225989754" MODIFIED="1408231258328" TEXT="All CSS, JS, HTML files are compiled into"/>
</node>
</node>
</map>
