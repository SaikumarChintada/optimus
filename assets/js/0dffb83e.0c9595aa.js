"use strict";(self.webpackChunkoptimus=self.webpackChunkoptimus||[]).push([[5075],{3905:function(e,t,r){r.d(t,{Zo:function(){return p},kt:function(){return m}});var n=r(7294);function i(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function o(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function a(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?o(Object(r),!0).forEach((function(t){i(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):o(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function s(e,t){if(null==e)return{};var r,n,i=function(e,t){if(null==e)return{};var r,n,i={},o=Object.keys(e);for(n=0;n<o.length;n++)r=o[n],t.indexOf(r)>=0||(i[r]=e[r]);return i}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(n=0;n<o.length;n++)r=o[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(i[r]=e[r])}return i}var u=n.createContext({}),l=function(e){var t=n.useContext(u),r=t;return e&&(r="function"==typeof e?e(t):a(a({},t),e)),r},p=function(e){var t=l(e.components);return n.createElement(u.Provider,{value:t},e.children)},c={inlineCode:"code",wrapper:function(e){var t=e.children;return n.createElement(n.Fragment,{},t)}},d=n.forwardRef((function(e,t){var r=e.components,i=e.mdxType,o=e.originalType,u=e.parentName,p=s(e,["components","mdxType","originalType","parentName"]),d=l(r),m=i,h=d["".concat(u,".").concat(m)]||d[m]||c[m]||o;return r?n.createElement(h,a(a({ref:t},p),{},{components:r})):n.createElement(h,a({ref:t},p))}));function m(e,t){var r=arguments,i=t&&t.mdxType;if("string"==typeof e||i){var o=r.length,a=new Array(o);a[0]=d;var s={};for(var u in t)hasOwnProperty.call(t,u)&&(s[u]=t[u]);s.originalType=e,s.mdxType="string"==typeof e?e:i,a[1]=s;for(var l=2;l<o;l++)a[l]=r[l];return n.createElement.apply(null,a)}return n.createElement.apply(null,r)}d.displayName="MDXCreateElement"},3208:function(e,t,r){r.r(t),r.d(t,{frontMatter:function(){return s},contentTitle:function(){return u},metadata:function(){return l},toc:function(){return p},default:function(){return d}});var n=r(7462),i=r(3366),o=(r(7294),r(3905)),a=["components"],s={id:"roadmap",title:"Roadmap",slug:"/roadmap"},u=void 0,l={unversionedId:"roadmap",id:"roadmap",isDocsHomePage:!1,title:"Roadmap",description:"This is a live document which gives an idea for all the users of Optimus on what are we upto.",source:"@site/docs/roadmap.md",sourceDirName:".",slug:"/roadmap",permalink:"/optimus/docs/roadmap",editUrl:"https://github.com/odpf/optimus/edit/master/docs/docs/roadmap.md",tags:[],version:"current",lastUpdatedBy:"Yash Bhardwaj",lastUpdatedAt:1649677473,formattedLastUpdatedAt:"4/11/2022",frontMatter:{id:"roadmap",title:"Roadmap",slug:"/roadmap"},sidebar:"docsSidebar",previous:{title:"Contributing",permalink:"/optimus/docs/contribute/contributing"},next:{title:"API",permalink:"/optimus/docs/reference/api"}},p=[{value:"Secret Management in Optimus",id:"secret-management-in-optimus",children:[]},{value:"Test Your Jobs",id:"test-your-jobs",children:[]},{value:"Telemetry",id:"telemetry",children:[]},{value:"Add More Plugins",id:"add-more-plugins",children:[]},{value:"Task Versioning",id:"task-versioning",children:[]},{value:"Improved Window Support",id:"improved-window-support",children:[]},{value:"SLA Tracking",id:"sla-tracking",children:[]},{value:"SQLite Support",id:"sqlite-support",children:[]},{value:"Custom Macro Support",id:"custom-macro-support",children:[]},{value:"Inbuilt Testing F/w",id:"inbuilt-testing-fw",children:[]},{value:"Inter Task/Hook Communication",id:"inter-taskhook-communication",children:[]}],c={toc:p};function d(e){var t=e.components,r=(0,i.Z)(e,a);return(0,o.kt)("wrapper",(0,n.Z)({},c,r,{components:t,mdxType:"MDXLayout"}),(0,o.kt)("p",null,"This is a live document which gives an idea for all the users of Optimus on what are we upto. "),(0,o.kt)("h3",{id:"secret-management-in-optimus"},"Secret Management in Optimus"),(0,o.kt)("p",null,"As optimus meant to deal with various warehouse systems & with the plugin support provides the capability to interact with other third party systems.\nIt brings a need for proper secret management to store & use securely for all users onboard."),(0,o.kt)("h3",{id:"test-your-jobs"},"Test Your Jobs"),(0,o.kt)("p",null,"Giving a provision for users to test the jobs before deploying helps users with faster feedbacks."),(0,o.kt)("h3",{id:"telemetry"},"Telemetry"),(0,o.kt)("p",null,"With proper monitoring we can get many insights into Optimus, which helps in debugging any failures. May not be a direct end user feature but this is very important.  "),(0,o.kt)("h3",{id:"add-more-plugins"},"Add More Plugins"),(0,o.kt)("p",null,"Once the data is analyzed in warehouse, there is always a need for getting the data out of the system for visualizations or for consumption. This is a constant effort to improve the ecosystem that optimus supports.\nThe plugins that we will be adding support is to pull data from BQ to Kafka, JDBC, FTP."),(0,o.kt)("h3",{id:"task-versioning"},"Task Versioning"),(0,o.kt)("p",null,"Versioning of tasks comes to handy when there is a time significance associated to task.\nOn replay, an older version of task has to run which was active at that time and the newer version on the coming days."),(0,o.kt)("h3",{id:"improved-window-support"},"Improved Window Support"),(0,o.kt)("p",null,"The current windowing which is used by for automated dependency resoultion & the macros which are derived from it are being used for input data filtering is little confusing & limiting in nature.\nThis will be an effort to easy the same."),(0,o.kt)("h3",{id:"sla-tracking"},"SLA Tracking"),(0,o.kt)("p",null,"Giving a provision for defining the SLAs & providing a dashboard to visualize how the slas are met is a must.\nWith this users will be able to monitor any slas that are breached out of the box."),(0,o.kt)("h3",{id:"sqlite-support"},"SQLite Support"),(0,o.kt)("p",null,"With support of SQLlite database, just helps users to kick start Optimus fast & easy to try on without having a dependency on postgres."),(0,o.kt)("h3",{id:"custom-macro-support"},"Custom Macro Support"),(0,o.kt)("p",null,"Custom Macros will unleash many capabilities, this will help users to template their queries to avoid any duplication."),(0,o.kt)("h3",{id:"inbuilt-testing-fw"},"Inbuilt Testing F/w"),(0,o.kt)("p",null,"Currently Optimus relies on Predator for Quality Checks, instead of relying on predator which is not extensible & supports only BQ, providing a capability to test the job runs directly."),(0,o.kt)("h3",{id:"inter-taskhook-communication"},"Inter Task/Hook Communication"),(0,o.kt)("p",null,"As we scale, there are situations where tasks & hooks has to share information directly instead of relying on another system. This opens up many capabilities."))}d.isMDXComponent=!0}}]);