
		var form_textbox = Vue.extend({
			template: '\
			<input type="text" />	',
		props: [],
	created(){
		console.log('created');
	},
	mounted(){
     // alert('a');
   }})
		Vue.component('form_textbox', form_textbox)
		
		var form_textarea = Vue.extend({
			template: '\
			<textarea rows="3" cols="20">\
 这是一个textarea\
</textarea>	',
		})
		Vue.component('form_textarea', form_textarea)
		
		var navmenu = Vue.extend({
			template: '\
			<div class="flex fixnav">\
	<router-link v-bind:to="n.path" v-for="n,index in list" v-bind:key="index">\
		<i v-bind:class="n.icon" class="iconfont"></i><p v-text="n.text">-</p>\
	</router-link>\
</div>	',
		data(){
return{
	list:[
      {
          "path":"/",
          "img":"/js/easyui/themes/icons/8.png",
		  "icon":"icon-home",
          "text":"首页"
      },
	  {
          "path":"/msg",
          "img":"/js/easyui/themes/icons/44.png",
		  "icon":"icon-icon--",
          "text":"应用"
      },
	  {
          "path":"/kecheng",
          "img":"/js/easyui/themes/icons/10.png",
		  "icon":"icon-fangdajing",
          "text":"组件"
      },
      {
          "path":"/center",
          "img":"/js/easyui/themes/icons/43.png",
		  "icon":"icon-gerenzhongxinwode",
          "text":"我的"
      }
	]
}
},
created(){
		console.log('nav created');
},
mounted(){
        console.log('nav mounted');
}})
		Vue.component('navmenu', navmenu)
		
		var app_notfound = Vue.extend({
			template: '\
			<div>\
            error:404\
        </div>	',
		})
		Vue.component('app_notfound', app_notfound)
		
		var weui_msg = Vue.extend({
			template: '\
			<div class="weui-msg" style="padding-top:0px;">\
  <div class="weui-msg__icon-area" style="margin-bottom:2px;">\
    <img style="width:100%;" src="/images/pic1.jpg" />\
  </div>\
\
</div>	',
		})
Vue.component('weui_msg', weui_msg)

 

var navmodule = Vue.extend({
			template: '\
	<div class="weui-grids grids-small" style="">\
      <a href="" class="weui-grid js_grid">\
        <div class="weui-grid__icon">\
          <img src="/images/icon/icon2.jpg" class="gray" alt="">\
        </div>\
        <p class="weui-grid__label gray">\
          课程订购\
        </p>\
      </a>\
      <a href="" class="weui-grid js_grid">\
        <div class="weui-grid__icon">\
          <img src="/images/icon/icon4.jpg" class="gray"  alt="">\
        </div>\
        <p class="weui-grid__label gray">\
          上课记录\
        </p>\
      </a>\
      <a href="" class="weui-grid js_grid">\
        <div class="weui-grid__icon">\
          <img src="/images/icon/icon3.jpg" class="gray"  alt="">\
        </div>\
        <p class="weui-grid__label gray">\
          视频课程\
        </p>\
      </a>\
      <a href="" class="weui-grid js_grid">\
        <div class="weui-grid__icon">\
          <img src="/images/icon/icon1.jpg" class="gray"  alt="">\
        </div>\
        <p class="weui-grid__label gray">\
          精品试卷\
        </p>\
      </a>\
		<a href="" class="weui-grid js_grid">\
        <div class="weui-grid__icon">\
          <img src="/images/icon/icon5.jpg" class="gray"  alt="">\
        </div>\
        <p class="weui-grid__label gray">\
          交费记录\
        </p>\
      </a>\
		<a href="" class="weui-grid js_grid">\
        <div class="weui-grid__icon">\
          <img src="/images/icon/icon8.jpg" class="gray"  alt="">\
        </div>\
        <p class="weui-grid__label gray">\
          幸运抽奖\
        </p>\
      </a>\
		<a href="" class="weui-grid js_grid">\
        <div class="weui-grid__icon">\
          <img src="/images/icon/icon7.jpg" class="gray"  alt="">\
        </div>\
        <p class="weui-grid__label gray">\
          积分兑换\
        </p>\
      </a>\
		<a href="" class="weui-grid js_grid">\
        <div class="weui-grid__icon">\
          <img src="/images/icon/icon5.jpg" class="gray"  alt="">\
        </div>\
        <p class="weui-grid__label gray">\
          上课视频\
        </p>\
      </a>\
    </div>\
	',
		})
Vue.component('navmodule', navmodule)









var xxms_topmsg = Vue.extend({
		template: '\
<div data-v-dec7aa6c="" class="bgwhite clear">\
   <div data-v-dec7aa6c="" class="datashows">\
    <p data-v-dec7aa6c="">账户余额（元）</p> \
    <h1 data-v-dec7aa6c="" class="price_cls">890.00</h1> \
    <div data-v-dec7aa6c="" class="flexbtn datanum_cls">\
     <p data-v-dec7aa6c="">已上课时<span data-v-dec7aa6c="">8</span></p> \
     <p data-v-dec7aa6c="">我的课表<span data-v-dec7aa6c="">12</span></p> \
     <p data-v-dec7aa6c="">当前积分<span data-v-dec7aa6c="">189</span></p>\
    </div>\
   </div>\
  </div>\
		',
		})
Vue.component('xxms_topmsg', xxms_topmsg)




var xxms_building = Vue.extend({
			template: '\
			<div class="weui-msg" style="padding-top:0px;">\
  <div class="weui-msg__icon-area" style="margin-bottom:2px;">\
    <img style="width:100%;" src="/images/building.jpg" />\
				<p data-v-7eb78c68="" class="tips" style="font-size:21px;">敬请期待</p>\
				<p data-v-7eb78c68="" class="tips">功能正在开发中...</p>\
  </div>\
\
</div>	',
		})
Vue.component('xxms_building', xxms_building)




var xxms_my = Vue.extend({
	template: '\
	<div data-v-ed66b0d2="" class="tit_cls"><i data-v-ed66b0d2="" class="iconfont icon-left"></i>我的账号\
	</div>\
	',
		})
Vue.component('xxms_my', xxms_my)











		