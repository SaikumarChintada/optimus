"use strict";(self.webpackChunkoptimus=self.webpackChunkoptimus||[]).push([[4212],{3905:function(t,e,n){n.d(e,{Zo:function(){return s},kt:function(){return c}});var a=n(7294);function r(t,e,n){return e in t?Object.defineProperty(t,e,{value:n,enumerable:!0,configurable:!0,writable:!0}):t[e]=n,t}function o(t,e){var n=Object.keys(t);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(t);e&&(a=a.filter((function(e){return Object.getOwnPropertyDescriptor(t,e).enumerable}))),n.push.apply(n,a)}return n}function i(t){for(var e=1;e<arguments.length;e++){var n=null!=arguments[e]?arguments[e]:{};e%2?o(Object(n),!0).forEach((function(e){r(t,e,n[e])})):Object.getOwnPropertyDescriptors?Object.defineProperties(t,Object.getOwnPropertyDescriptors(n)):o(Object(n)).forEach((function(e){Object.defineProperty(t,e,Object.getOwnPropertyDescriptor(n,e))}))}return t}function l(t,e){if(null==t)return{};var n,a,r=function(t,e){if(null==t)return{};var n,a,r={},o=Object.keys(t);for(a=0;a<o.length;a++)n=o[a],e.indexOf(n)>=0||(r[n]=t[n]);return r}(t,e);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(t);for(a=0;a<o.length;a++)n=o[a],e.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(t,n)&&(r[n]=t[n])}return r}var p=a.createContext({}),u=function(t){var e=a.useContext(p),n=e;return t&&(n="function"==typeof t?t(e):i(i({},e),t)),n},s=function(t){var e=u(t.components);return a.createElement(p.Provider,{value:e},t.children)},d={inlineCode:"code",wrapper:function(t){var e=t.children;return a.createElement(a.Fragment,{},e)}},m=a.forwardRef((function(t,e){var n=t.components,r=t.mdxType,o=t.originalType,p=t.parentName,s=l(t,["components","mdxType","originalType","parentName"]),m=u(n),c=r,g=m["".concat(p,".").concat(c)]||m[c]||d[c]||o;return n?a.createElement(g,i(i({ref:e},s),{},{components:n})):a.createElement(g,i({ref:e},s))}));function c(t,e){var n=arguments,r=e&&e.mdxType;if("string"==typeof t||r){var o=n.length,i=new Array(o);i[0]=m;var l={};for(var p in e)hasOwnProperty.call(e,p)&&(l[p]=e[p]);l.originalType=t,l.mdxType="string"==typeof t?t:r,i[1]=l;for(var u=2;u<o;u++)i[u]=n[u];return a.createElement.apply(null,i)}return a.createElement.apply(null,n)}m.displayName="MDXCreateElement"},1406:function(t,e,n){n.r(e),n.d(e,{frontMatter:function(){return l},contentTitle:function(){return p},metadata:function(){return u},toc:function(){return s},default:function(){return m}});var a=n(7462),r=n(3366),o=(n(7294),n(3905)),i=["components"],l={id:"publishing-bigquery-to-kafka",title:"Publishing from bigQuery to kafka"},p=void 0,u={unversionedId:"guides/publishing-bigquery-to-kafka",id:"guides/publishing-bigquery-to-kafka",isDocsHomePage:!1,title:"Publishing from bigQuery to kafka",description:"Following are the steps to publish BQ data to Kafka:",source:"@site/docs/guides/publishing-from-bigquery-to-kafka.md",sourceDirName:"guides",slug:"/guides/publishing-bigquery-to-kafka",permalink:"/optimus/docs/guides/publishing-bigquery-to-kafka",editUrl:"https://github.com/odpf/optimus/edit/master/docs/docs/guides/publishing-from-bigquery-to-kafka.md",tags:[],version:"current",lastUpdatedBy:"sravankorumilli",lastUpdatedAt:1650468189,formattedLastUpdatedAt:"4/20/2022",frontMatter:{id:"publishing-bigquery-to-kafka",title:"Publishing from bigQuery to kafka"}},s=[{value:"Config:",id:"config",children:[]},{value:"Bigquery and Protobuf Data Type Mapping",id:"bigquery-and-protobuf-data-type-mapping",children:[]}],d={toc:s};function m(t){var e=t.components,n=(0,r.Z)(t,i);return(0,o.kt)("wrapper",(0,a.Z)({},d,n,{components:e,mdxType:"MDXLayout"}),(0,o.kt)("p",null,"Following are the steps to publish BQ data to Kafka:"),(0,o.kt)("ol",null,(0,o.kt)("li",{parentName:"ol"},"Raise a ticket with DE to create the proto, soon this will be automated."),(0,o.kt)("li",{parentName:"ol"},"As part of Job Specification answer the questions related to pushing data to Kafka."),(0,o.kt)("li",{parentName:"ol"},"Commit and Push")),(0,o.kt)("p",null,"Publishing from Bigquery to Kafka is done using an Optimus hook called ",(0,o.kt)("inlineCode",{parentName:"p"},"Transporter"),".\nOnce data is pushed to Kafka you can use firehose to consume messages for your needs.\nIn order to add hook to an existing Job, run the following command and answer the\ncorresponding prompts:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre"},"$ ./optimus job addhook\n? Select a Job example_job\n? Which hook to run? transporter\n? Filter expression for extracting transformation rows? event_timestamp >= '{{.DSTART}}' AND event_timestamp < '{{.DEND}}'\n")),(0,o.kt)("p",null,"Note: this is not available for public use at the moment"),(0,o.kt)("h3",{id:"config"},"Config:"),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},"Filter expression: Expression is used as a where clause to restrict the number of rows to only push the ones\nthat are affected by current transformation. Expression can be templated with DSTART and DEND macros which will be replaced with the window for which the current transformation is getting executed similar to what we do in Job specifications.")),(0,o.kt)("p",null,"After this, existing job.yaml file will get updated with the new hooks config and the job specification would look like below:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-yaml"},"version: 1\nname: example_job\nowner: example@example.com\nschedule:\n  start_date: \"2021-02-18\"\n  interval: 0 3 * * *\nbehavior:\n  depends_on_past: false\n  catch_up: true\ntask:\n  name: bq2bq\n  config:\n    DATASET: data\n    LOAD_METHOD: APPEND\n    PROJECT: example\n    SQL_TYPE: STANDARD\n    TABLE: hello_table\n  window:\n    size: 24h\n    offset: \"0\"\n    truncate_to: d\ndependencies: []\nhooks:\n- name: transporter\n  config:\n    BQ_DATASET: '{{.TASK__DATASET}}' # inherited from task configs\n    BQ_PROJECT: '{{.TASK__PROJECT}}'\n    BQ_TABLE: '{{.TASK__TABLE}}'\n    FILTER_EXPRESSION: 'event_timestamp >= \"{{.DSTART}}\" AND event_timestamp < \"{{.DEND}}\"'\n    KAFKA_TOPIC: optimus_example-data-hello_table\n    PRODUCER_CONFIG_BOOTSTRAP_SERVERS: '{{.GLOBAL__TRANSPORTER_KAFKA_BROKERS}}'\n    PROTO_SCHEMA: example.data.HelloTable\n    STENCIL_URL: '{{.GLOBAL__TRANSPORTER_KAFKA_BROKERS}}' # will be defined as global config\n")),(0,o.kt)("p",null,"Some configuration would be auto-generated by Optimus.\nYou can now commit and push the changes in Job specification."),(0,o.kt)("p",null,"There are standards that we want users to maintain but haven\u2019t enforced:"),(0,o.kt)("ol",null,(0,o.kt)("li",{parentName:"ol"},"Single table to topic mapping to be maintained & following a naming convention\nhelps in better discoverability. Topic names are auto populated with the same naming convention."),(0,o.kt)("li",{parentName:"ol"},"Kafka Brokers are not modifiable & pre-selected per every entity."),(0,o.kt)("li",{parentName:"ol"},"All Protos & Stencil are managed & doesn\u2019t need any user intervention.")),(0,o.kt)("h3",{id:"bigquery-and-protobuf-data-type-mapping"},"Bigquery and Protobuf Data Type Mapping"),(0,o.kt)("table",null,(0,o.kt)("thead",{parentName:"table"},(0,o.kt)("tr",{parentName:"thead"},(0,o.kt)("th",{parentName:"tr",align:null},"Bigquery Type"),(0,o.kt)("th",{parentName:"tr",align:null},"Protobuf Type"))),(0,o.kt)("tbody",{parentName:"table"},(0,o.kt)("tr",{parentName:"tbody"},(0,o.kt)("td",{parentName:"tr",align:null},"STRING"),(0,o.kt)("td",{parentName:"tr",align:null},"string")),(0,o.kt)("tr",{parentName:"tbody"},(0,o.kt)("td",{parentName:"tr",align:null},"BYTES"),(0,o.kt)("td",{parentName:"tr",align:null},"bytes")),(0,o.kt)("tr",{parentName:"tbody"},(0,o.kt)("td",{parentName:"tr",align:null},"INTEGER"),(0,o.kt)("td",{parentName:"tr",align:null},"int64")),(0,o.kt)("tr",{parentName:"tbody"},(0,o.kt)("td",{parentName:"tr",align:null},"FLOAT"),(0,o.kt)("td",{parentName:"tr",align:null},"float")),(0,o.kt)("tr",{parentName:"tbody"},(0,o.kt)("td",{parentName:"tr",align:null},"BOOLEAN"),(0,o.kt)("td",{parentName:"tr",align:null},"bool")),(0,o.kt)("tr",{parentName:"tbody"},(0,o.kt)("td",{parentName:"tr",align:null},"TIMESTAMP"),(0,o.kt)("td",{parentName:"tr",align:null},"google.protobuf.Timestamp")),(0,o.kt)("tr",{parentName:"tbody"},(0,o.kt)("td",{parentName:"tr",align:null},"DATE"),(0,o.kt)("td",{parentName:"tr",align:null},"string")),(0,o.kt)("tr",{parentName:"tbody"},(0,o.kt)("td",{parentName:"tr",align:null},"TIME"),(0,o.kt)("td",{parentName:"tr",align:null},"string")),(0,o.kt)("tr",{parentName:"tbody"},(0,o.kt)("td",{parentName:"tr",align:null},"DATETIME"),(0,o.kt)("td",{parentName:"tr",align:null},"string")),(0,o.kt)("tr",{parentName:"tbody"},(0,o.kt)("td",{parentName:"tr",align:null},"NUMERIC"),(0,o.kt)("td",{parentName:"tr",align:null},"float")),(0,o.kt)("tr",{parentName:"tbody"},(0,o.kt)("td",{parentName:"tr",align:null},"GEOGRAPHY"),(0,o.kt)("td",{parentName:"tr",align:null},"string")))))}m.isMDXComponent=!0}}]);