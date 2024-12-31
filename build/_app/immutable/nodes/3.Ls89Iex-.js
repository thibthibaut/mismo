import{a as h,t as b}from"../chunks/disclose-version.CTQXdc_7.js";import"../chunks/legacy.CBGdgsew.js";import{p as F,l as H,a as O,b as P,g as t,c,f as y,s as l,r as v,t as S,m as n,d as p}from"../chunks/runtime.BYmxHEYc.js";import{e as k,s as W}from"../chunks/render.B-6WVgCH.js";import{i as C}from"../chunks/props.VDm8niNF.js";import{r as G,b as z,W as J,t as R,s as U,f as q}from"../chunks/index.Cn-1KTV1.js";import{i as A}from"../chunks/lifecycle.QqYsvlCh.js";var K=b(`<input name="name" placeholder="Ton Blaze" class="px-6 py-3
                   bg-white
                   border-2 border-violet-200
                   focus:border-violet-400 focus:ring-2 focus:ring-violet-200
                   rounded-lg
                   shadow-sm
                   placeholder-violet-300
                   text-violet-600
                   min-w-[200px]
                   transition-all duration-200
                   outline-none"> <button class="px-6 py-3 
                   bg-gradient-to-r from-violet-600 to-purple-600 
                   hover:from-violet-700 hover:to-purple-700
                   text-white font-medium rounded-lg 
                   shadow-lg hover:shadow-xl
                   transition-all duration-200 
                   min-w-[200px]">Créer la partie</button>`,1),L=b('<div class="mt-4 p-4 bg-violet-50 rounded-lg border-2 border-violet-200"><p class="text-violet-600 mb-2"> </p> <div class="flex items-center space-x-2"><input readonly="" class="px-4 py-2 bg-white rounded border border-violet-200 text-violet-600"> <button class="p-2 bg-violet-100 hover:bg-violet-200 rounded-lg transition-colors"><svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-violet-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3"></path></svg></button></div></div> <!>',1),Q=b('<div class="flex flex-col space-y-4 items-center"><!> <!></div>');function ae(T,j){F(j,!1);let o=n(""),i=n(""),m=n(""),w=n(!1),d=n(!0);async function I(){if(!t(o).trim()){alert("Please enter your name first");return}try{const e=await fetch("/create-game",{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify({playerName:t(o)})});if(!e.ok){const r=await e.text();throw new Error(r||"Failed to create game.")}const a=await e.json();p(i,a.gameId),p(m,`${window.location.origin}/join/${t(i)}`),p(d,!1),p(w,!0)}catch(e){console.error("Create Game Error:",e),alert(`Error: ${e.message}`)}}async function N(){try{await navigator.clipboard.writeText(t(m))}catch(e){console.error("Failed to copy text: ",e)}}H(()=>t(d),()=>{console.log("showCreateGame:",t(d))}),O(),A();var f=Q(),x=c(f);{var $=e=>{var a=K(),r=y(a);G(r);var s=l(r,2);z(r,()=>t(o),u=>p(o,u)),k("click",s,I),h(e,a)};C(x,e=>{t(d)&&e($)})}var E=l(x,2);{var M=e=>{var a=L(),r=y(a),s=c(r),u=c(s);v(s);var _=l(s,2),g=c(_);G(g);var B=l(g,2);v(_),v(r);var D=l(r,2);J(D,{get gameID(){return t(i)},get playerName(){return t(o)}}),S(()=>{W(u,`Game ID: ${t(i)??""}`),U(g,t(m))}),k("click",B,N),R(3,r,()=>q),h(e,a)};C(E,e=>{t(w)&&t(i)&&t(o)&&e(M)})}v(f),h(T,f),P()}export{ae as component};
