import{S as N,C as V,D as z,F as R,G as J,d as h,H as D,U as c,g as P,I as F,J as Q,K as W,L as X,w as k,A as C,M as p,E as ee,N as ae,O as re,P as te,Q as q,R as M,x as U,T as G,B as ne,V as ie,W as fe,X as se,Y as ue,i as H,Z as le,_ as _e,$ as K,a0 as ve,v as de,a1 as ce,a2 as oe,a3 as be,o as Z,a4 as ge,a5 as ye,a6 as he}from"./runtime.BYmxHEYc.js";function w(n,u=null,g){if(typeof n!="object"||n===null||N in n)return n;const v=W(n);if(v!==V&&v!==z)return n;var i=new Map,_=X(n),o=R(0);_&&i.set("length",R(n.length));var y;return new Proxy(n,{defineProperty(f,e,a){(!("value"in a)||a.configurable===!1||a.enumerable===!1||a.writable===!1)&&J();var t=i.get(e);return t===void 0?(t=R(a.value),i.set(e,t)):h(t,w(a.value,y)),!0},deleteProperty(f,e){var a=i.get(e);if(a===void 0)e in f&&i.set(e,R(c));else{if(_&&typeof e=="string"){var t=i.get("length"),r=Number(e);Number.isInteger(r)&&r<t.v&&h(t,r)}h(a,c),$(o)}return!0},get(f,e,a){var d;if(e===N)return n;var t=i.get(e),r=e in f;if(t===void 0&&(!r||(d=D(f,e))!=null&&d.writable)&&(t=R(w(r?f[e]:c,y)),i.set(e,t)),t!==void 0){var s=P(t);return s===c?void 0:s}return Reflect.get(f,e,a)},getOwnPropertyDescriptor(f,e){var a=Reflect.getOwnPropertyDescriptor(f,e);if(a&&"value"in a){var t=i.get(e);t&&(a.value=P(t))}else if(a===void 0){var r=i.get(e),s=r==null?void 0:r.v;if(r!==void 0&&s!==c)return{enumerable:!0,configurable:!0,value:s,writable:!0}}return a},has(f,e){var s;if(e===N)return!0;var a=i.get(e),t=a!==void 0&&a.v!==c||Reflect.has(f,e);if(a!==void 0||F!==null&&(!t||(s=D(f,e))!=null&&s.writable)){a===void 0&&(a=R(t?w(f[e],y):c),i.set(e,a));var r=P(a);if(r===c)return!1}return t},set(f,e,a,t){var E;var r=i.get(e),s=e in f;if(_&&e==="length")for(var d=a;d<r.v;d+=1){var m=i.get(d+"");m!==void 0?h(m,c):d in f&&(m=R(c),i.set(d+"",m))}r===void 0?(!s||(E=D(f,e))!=null&&E.writable)&&(r=R(void 0),h(r,w(a,y)),i.set(e,r)):(s=r.v!==c,h(r,w(a,y)));var b=Reflect.getOwnPropertyDescriptor(f,e);if(b!=null&&b.set&&b.set.call(t,a),!s){if(_&&typeof e=="string"){var S=i.get("length"),O=Number(e);Number.isInteger(O)&&O>=S.v&&h(S,O+1)}$(o)}return!0},ownKeys(f){P(o);var e=Reflect.ownKeys(f).filter(r=>{var s=i.get(r);return s===void 0||s.v!==c});for(var[a,t]of i)t.v!==c&&!(a in f)&&e.push(a);return e},setPrototypeOf(){Q()}})}function $(n,u=1){h(n,n.v+u)}function me(n,u,g=!1){C&&p();var v=n,i=null,_=null,o=c,y=g?ee:0,f=!1;const e=(t,r=!0)=>{f=!0,a(r,t)},a=(t,r)=>{if(o===(o=t))return;let s=!1;if(C){const d=v.data===ae;!!o===d&&(v=re(),te(v),q(!1),s=!0)}o?(i?M(i):r&&(i=U(()=>r(v))),_&&G(_,()=>{_=null})):(_?M(_):r&&(_=U(()=>r(v))),i&&G(i,()=>{i=null})),s&&q(!0)};k(()=>{f=!1,u(e),f||a(null,null)},y),C&&(v=ne)}let A=!1;function Pe(n){var u=A;try{return A=!1,[n(),A]}finally{A=u}}function j(n){for(var u=F,g=F;u!==null&&!(u.f&(le|_e));)u=u.parent;try{return K(u),n()}finally{K(g)}}function Ee(n,u,g,v){var Y;var i=(g&ve)!==0,_=!de||(g&ce)!==0,o=(g&oe)!==0,y=(g&he)!==0,f=!1,e;o?[e,f]=Pe(()=>n[u]):e=n[u];var a=N in n||be in n,t=((Y=D(n,u))==null?void 0:Y.set)??(a&&o&&u in n?l=>n[u]=l:void 0),r=v,s=!0,d=!1,m=()=>(d=!0,s&&(s=!1,y?r=H(v):r=v),r);e===void 0&&v!==void 0&&(t&&_&&ie(),e=m(),t&&t(e));var b;if(_)b=()=>{var l=n[u];return l===void 0?m():(s=!0,d=!1,l)};else{var S=j(()=>(i?Z:ge)(()=>n[u]));S.f|=fe,b=()=>{var l=P(S);return l!==void 0&&(r=void 0),l===void 0?r:l}}if(!(g&se))return b;if(t){var O=n.$$legacy;return function(l,I){return arguments.length>0?((!_||!I||O||f)&&t(I?b():l),l):b()}}var E=!1,B=!1,L=ye(e),T=j(()=>Z(()=>{var l=b(),I=P(L);return E?(E=!1,B=!0,I):(B=!1,L.v=l)}));return i||(T.equals=ue),function(l,I){if(arguments.length>0){const x=I?P(T):_&&o?w(l):l;return T.equals(x)||(E=!0,h(L,x),d&&r!==void 0&&(r=x),H(()=>P(T))),l}return P(T)}}export{w as a,me as i,Ee as p};
