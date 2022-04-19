"use strict";(self.webpackChunkoptimus=self.webpackChunkoptimus||[]).push([[3436],{3905:function(e,t,n){n.d(t,{Zo:function(){return u},kt:function(){return d}});var a=n(7294);function i(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function r(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);t&&(a=a.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,a)}return n}function o(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?r(Object(n),!0).forEach((function(t){i(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):r(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function l(e,t){if(null==e)return{};var n,a,i=function(e,t){if(null==e)return{};var n,a,i={},r=Object.keys(e);for(a=0;a<r.length;a++)n=r[a],t.indexOf(n)>=0||(i[n]=e[n]);return i}(e,t);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);for(a=0;a<r.length;a++)n=r[a],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(i[n]=e[n])}return i}var p=a.createContext({}),s=function(e){var t=a.useContext(p),n=t;return e&&(n="function"==typeof e?e(t):o(o({},t),e)),n},u=function(e){var t=s(e.components);return a.createElement(p.Provider,{value:t},e.children)},c={inlineCode:"code",wrapper:function(e){var t=e.children;return a.createElement(a.Fragment,{},t)}},m=a.forwardRef((function(e,t){var n=e.components,i=e.mdxType,r=e.originalType,p=e.parentName,u=l(e,["components","mdxType","originalType","parentName"]),m=s(n),d=i,f=m["".concat(p,".").concat(d)]||m[d]||c[d]||r;return n?a.createElement(f,o(o({ref:t},u),{},{components:n})):a.createElement(f,o({ref:t},u))}));function d(e,t){var n=arguments,i=t&&t.mdxType;if("string"==typeof e||i){var r=n.length,o=new Array(r);o[0]=m;var l={};for(var p in t)hasOwnProperty.call(t,p)&&(l[p]=t[p]);l.originalType=e,l.mdxType="string"==typeof e?e:i,o[1]=l;for(var s=2;s<r;s++)o[s]=n[s];return a.createElement.apply(null,o)}return a.createElement.apply(null,n)}m.displayName="MDXCreateElement"},2783:function(e,t,n){n.r(t),n.d(t,{frontMatter:function(){return l},contentTitle:function(){return p},metadata:function(){return s},toc:function(){return u},default:function(){return m}});var a=n(7462),i=n(3366),r=(n(7294),n(3905)),o=["components"],l={id:"configuration",title:"Configurations"},p=void 0,s={unversionedId:"getting-started/configuration",id:"getting-started/configuration",isDocsHomePage:!1,title:"Configurations",description:"Client Configuration",source:"@site/docs/getting-started/configuration.md",sourceDirName:"getting-started",slug:"/getting-started/configuration",permalink:"/optimus/docs/getting-started/configuration",editUrl:"https://github.com/odpf/optimus/edit/master/docs/docs/getting-started/configuration.md",tags:[],version:"current",lastUpdatedBy:"Anwar Hidayat",lastUpdatedAt:1650360802,formattedLastUpdatedAt:"4/19/2022",frontMatter:{id:"configuration",title:"Configurations"},sidebar:"docsSidebar",previous:{title:"Installation",permalink:"/optimus/docs/getting-started/installation"},next:{title:"Using Optimus to create a Job",permalink:"/optimus/docs/guides/create-job"}},u=[{value:"Client Configuration",id:"client-configuration",children:[{value:"Initialization",id:"initialization",children:[]},{value:"Usage",id:"usage",children:[]}]},{value:"Server Configuration",id:"server-configuration",children:[]}],c={toc:u};function m(e){var t=e.components,n=(0,i.Z)(e,o);return(0,r.kt)("wrapper",(0,a.Z)({},c,n,{components:t,mdxType:"MDXLayout"}),(0,r.kt)("h2",{id:"client-configuration"},"Client Configuration"),(0,r.kt)("h3",{id:"initialization"},"Initialization"),(0,r.kt)("p",null,"Client configuration can be initialized using Optimus.\nTo do so, use Optimus ",(0,r.kt)("inlineCode",{parentName:"p"},"init")," command.\nAn interactive questionaire will be presented, such as below:"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-zsh"},"$ optimus init\n\n? What is the Optimus service host? [? for help] localhost:9100\n? What is the Optimus project name? local-project\n? What is the namespace name? namespace1\n? What is the type of data store for this namespace? bigquery\n? Do you want to add another namespace? No\nClient config is initialized successfully\nIf you want to modify, go to [optimus.yaml]\n")),(0,r.kt)("p",null,"After running the ",(0,r.kt)("inlineCode",{parentName:"p"},"init")," command, Optimus client config will be configured.\nAlong with it is the directories for the chosen namespaces.\nFor example, since the above example only specifies ",(0,r.kt)("em",{parentName:"p"},"namespace1"),", then directory ",(0,r.kt)("em",{parentName:"p"},"namespace1")," will be created along with ",(0,r.kt)("inlineCode",{parentName:"p"},"optimus.yaml")," client config."),(0,r.kt)("p",null,"The user can optionally update the ",(0,r.kt)("inlineCode",{parentName:"p"},"optimus.yaml")," file if necessary.\nFor clearer example, see client configuration example on ",(0,r.kt)("inlineCode",{parentName:"p"},"optimus.sample.yaml")),(0,r.kt)("h3",{id:"usage"},"Usage"),(0,r.kt)("p",null,"Optimus client configuration can be loaded from file (use ",(0,r.kt)("inlineCode",{parentName:"p"},"--config")," flag), or ",(0,r.kt)("inlineCode",{parentName:"p"},"optimus.yaml")," file in current working directory where the Optimus command is executed."),(0,r.kt)("hr",null),(0,r.kt)("p",null,(0,r.kt)("strong",{parentName:"p"},"1. Using --config flag")),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-sh"},"$ optimus deploy --config /path/to/config/file.yaml\n")),(0,r.kt)("hr",null),(0,r.kt)("p",null,(0,r.kt)("strong",{parentName:"p"},"2. Using default optimus.yaml file")),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-sh"},"$ tree\n. # current project structure\n\u251c\u2500\u2500 namespace-1\n\u2502\xa0\xa0 \u2514\u2500\u2500 jobs\n\u2502\xa0\xa0 \u2514\u2500\u2500 resources\n\u251c\u2500\u2500 namespace-2\n\u2502\xa0\xa0 \u2514\u2500\u2500 jobs\n\u2502\xa0\xa0 \u2514\u2500\u2500 resources\n\u2514\u2500\u2500 optimus.yaml # use this file\n$ optimus deploy\n")),(0,r.kt)("hr",null),(0,r.kt)("p",null,"If both ",(0,r.kt)("inlineCode",{parentName:"p"},"--config")," flag and ",(0,r.kt)("inlineCode",{parentName:"p"},"optimus.yaml")," do exist, then the one will be used is the file defined in ",(0,r.kt)("inlineCode",{parentName:"p"},"--config")," flag."),(0,r.kt)("h2",{id:"server-configuration"},"Server Configuration"),(0,r.kt)("p",null,"See server configuration example on ",(0,r.kt)("inlineCode",{parentName:"p"},"config.sample.yaml")),(0,r.kt)("p",null,"Optimus server configuration can be loaded from file (use ",(0,r.kt)("inlineCode",{parentName:"p"},"--config")," flag), environment variable ",(0,r.kt)("inlineCode",{parentName:"p"},"OPTIMUS_<CONFIGNAME>"),", or ",(0,r.kt)("inlineCode",{parentName:"p"},"config.yaml")," file in executable directory."),(0,r.kt)("hr",null),(0,r.kt)("p",null,(0,r.kt)("strong",{parentName:"p"},"1. Using --config flag")),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-sh"},"$ optimus serve --config /path/to/config/file.yaml\n")),(0,r.kt)("hr",null),(0,r.kt)("p",null,(0,r.kt)("strong",{parentName:"p"},"2. Using environment variable")),(0,r.kt)("p",null,"All the configs can be passed as environment variables using ",(0,r.kt)("inlineCode",{parentName:"p"},"OPTIMUS_<CONFIG_NAME>")," convention. ",(0,r.kt)("inlineCode",{parentName:"p"},"<CONFIG_NAME>")," is the key name of config using ",(0,r.kt)("inlineCode",{parentName:"p"},"_")," as the path delimiter to concatenate between keys."),(0,r.kt)("p",null,"For example, to use environment variable, assuming the following configuration layout:"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-yaml"},"version: 1\nserve:\n  port: 9100\n  app_key: randomhash\n")),(0,r.kt)("p",null,"Here is the corresponding environment variable for the above"),(0,r.kt)("table",null,(0,r.kt)("thead",{parentName:"table"},(0,r.kt)("tr",{parentName:"thead"},(0,r.kt)("th",{parentName:"tr",align:null},"Configuration key"),(0,r.kt)("th",{parentName:"tr",align:null},"Environment variable"))),(0,r.kt)("tbody",{parentName:"table"},(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},"version"),(0,r.kt)("td",{parentName:"tr",align:null},"OPTIMUS_VERSION")),(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},"serve.port"),(0,r.kt)("td",{parentName:"tr",align:null},"OPTIMUS_PORT")),(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},"serve.app_key"),(0,r.kt)("td",{parentName:"tr",align:null},"OPTIMUS_SERVE_APP_KEY")))),(0,r.kt)("p",null,"Set the env variable using export"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-sh"},"$ export OPTIMUS_PORT=9100\n")),(0,r.kt)("hr",null),(0,r.kt)("p",null,(0,r.kt)("strong",{parentName:"p"},"3. Using default config.yaml from executable binary directory")),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-sh"},"$ which optimus\n/usr/local/bin/optimus\n")),(0,r.kt)("p",null,"So the ",(0,r.kt)("inlineCode",{parentName:"p"},"config.yaml")," file can be loaded on ",(0,r.kt)("inlineCode",{parentName:"p"},"/usr/local/bin/config.yaml")),(0,r.kt)("hr",null),(0,r.kt)("p",null,"If user specify configuration file using ",(0,r.kt)("inlineCode",{parentName:"p"},"--config")," flag, then any configs defined in env variable and default config.yaml from exec directory will not be loaded."),(0,r.kt)("p",null,"If user specify the env variable and user also have config.yaml from exec directory, then any key configs from env variable will override the key configs defined in config.yaml from exec directory."),(0,r.kt)("hr",null),(0,r.kt)("p",null,"App key is used to encrypt credentials and can be randomly generated using"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-shell"},"head -c 50 /dev/random | base64\n")),(0,r.kt)("p",null,"Just take the first 32 characters of the string."))}m.isMDXComponent=!0}}]);