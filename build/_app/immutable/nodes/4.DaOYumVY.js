import{a as d,t as c}from"../chunks/disclose-version.CTQXdc_7.js";import"../chunks/legacy.CBGdgsew.js";import{g as t,c as v,s as m,r as p,m as u,t as z,d as a}from"../chunks/runtime.BYmxHEYc.js";import{e as A,s as B}from"../chunks/render.B-6WVgCH.js";import{i as _}from"../chunks/props.VDm8niNF.js";import{r as y,b as D,t as I,W as C,f as P}from"../chunks/index.Cn-1KTV1.js";var E=c('<div class="text-red-500 bg-red-50 p-3 rounded-lg"> </div>'),F=c(`<div class="flex flex-col space-y-4 items-center"><!> <input name="gameID" placeholder="Game ID" class="px-6 py-3
                       bg-white
                       border-2 border-violet-200
                       focus:border-violet-400 focus:ring-2 focus:ring-violet-200
                       rounded-lg
                       shadow-sm
                       placeholder-violet-300
                       text-violet-600
                       min-w-[200px]
                       transition-all duration-200
                       outline-none"> <input name="name" placeholder="Ton Prénom" class="px-6 py-3
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
                       min-w-[200px]">Rejoindre la partie</button></div>`),H=c("<div><!></div>"),J=c('<div class="flex flex-col space-y-4 items-center"><!> <!></div>');function U(R){let i=u(""),n=u(""),h=u(!1),l=u("");function W(){if(!t(i).trim()){a(l,"Please enter your name");return}if(!t(n).trim()){a(l,"Please enter the game ID");return}a(h,!0)}var f=J(),w=v(f);{var j=r=>{var e=F(),s=v(e);{var N=o=>{var x=E(),q=v(x,!0);p(x),z(()=>B(q,t(l))),d(o,x)};_(s,o=>{t(l)&&o(N)})}var g=m(s,2);y(g);var b=m(g,2);y(b);var T=m(b,2);p(e),D(g,()=>t(n),o=>a(n,o)),D(b,()=>t(i),o=>a(i,o)),A("click",T,W),I(3,e,()=>P),d(r,e)};_(w,r=>{r(j)})}var k=m(w,2);{var G=r=>{var e=H(),s=v(e);C(s,{get gameID(){return t(n)},get playerName(){return t(i)}}),p(e),I(3,e,()=>P),d(r,e)};_(k,r=>{t(h)&&r(G)})}p(f),d(R,f)}export{U as component};
