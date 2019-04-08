var mainPlatform = {

	init: function(){

		this.bindEvent();
		// this.render(menu['home']);
	},

	bindEvent: function(){
		var self = this;
		// 顶部大菜单单击事件
		$(document).on('click', '.pf-nav-item', function() {
            $('.pf-nav-item').removeClass('current');
            $(this).addClass('current');

            // 渲染对应侧边菜单
            var m = $(this).data('menu');
            //self.render(menu[m]);
        });

        $(document).on('click', '.sider-nav li', function() {
            //$('.sider-nav li').removeClass('current');//合并菜单
            $(this).addClass('current');
            //$('iframe').attr('src', $(this).data('src'));
        });
		$(document).on('click', '.sider-nav-s li', function() {
            $('.sider-nav li').removeClass('active');
            $(this).addClass('active');
			/*
			if(!$('#tabs').tabs('exists',$(this).data('title'))){
				$('#tabs').tabs('add',{
					title:$(this).data('title'),
					content:'<iframe name="mainFrame" scrolling="auto" frameborder="0"  src="'+$(this).data('src')+'" style="width:100%;height:99%;"></iframe>',
					closable:true,
					width:$('#mainPanle').width()-10,
					height:$('#mainPanle').height()-26
				});
			}else{
				$('#tabs').tabs('select',$(this).data('title'));
			}
			*/
			$('#tabs').tabs('close',$(this).data('title'));
			if(!$('#tabs').tabs('exists',$(this).data('title'))){
				$('#tabs').tabs('add',{
					title:$(this).data('title'),
					content:'<iframe name="mainFrame" title="'+$(this).data('title')+'" scrolling="auto" frameborder="0"  src="'+$(this).data('src')+'" style="width:100%;height:99%;"></iframe>',
					closable:true,
					width:$('#mainPanle').width()-10,
					height:$('#mainPanle').height()-26
				});
			}else{
				$('#tabs').tabs('select',$(this).data('title'));
			}

        });

        $(document).on('click', '.pf-logout', function() {
			$.messager.confirm('确认', '<br/><br/>确定要退出系统吗?', function (r) {
				if (r) {
					location.href= '/adm/login'; 
				}
			});
			/*
            layer.confirm('您确定要退出吗？', {
              icon: 4,
			  title: '确定退出' //按钮
			}, function(){
			  location.href= 'login.html'; 
			});
			*/
        });
        //左侧菜单收起
        $(document).on('click', '.toggle-icon', function() {
            $(this).closest("#pf-bd").toggleClass("toggle");
            setTimeout(function(){
            	$(window).resize();
            },300)
        });

        $(document).on('click', '.pf-modify-pwd', function() {
            //$('#pf-page').find('iframe').eq(0).attr('src', 'backend/modify_pwd.html')
			addTab('修改密码','/adm/user/pwd');
        });

        $(document).on('click', '.pf-opt-name', function() {
            //$('#pf-page').find('iframe').eq(0).attr('src', 'backend/notice.html')
			addTab('用户信息','/adm/user/uinfo');
        });
	},
	
	render: function(menu){
		var current,
			html = ['<h2 class="pf-model-name"><span class="pf-sider-icon"></span><span class="pf-name">'+ menu.title +'</span></h2>'];

		html.push('<ul class="sider-nav">');
		for(var i = 0, len = menu.menu.length; i < len; i++){
			if(menu.menu[i].isCurrent){
				current = menu.menu[i];
				html.push('<li class="current" title="'+ menu.menu[i].title +'" data-src="'+ menu.menu[i].href +'"><a href="javascript:;"><img src="'+ menu.menu[i].icon +'"><span class="sider-nav-title">'+ menu.menu[i].title +'</span><i class="iconfont"></i></a></li>');
			}else{
				html.push('<li data-src="'+ menu.menu[i].href +'" title="'+ menu.menu[i].title +'"><a href="javascript:;"><img src="'+ menu.menu[i].icon +'"><span class="sider-nav-title">'+ menu.menu[i].title +'</span><i class="iconfont"></i></a></li>');
			}
		}
		html.push('</ul>');

		$('iframe').attr('src', current.href);
		$('#pf-sider').html(html.join(''));
	}

};

mainPlatform.init();
$.extend($.messager.defaults,{
		ok:"确定",
		cancel:"取消"
});
function addTab(title,url){
	//$('#tabs').tabs('close',title);
	if(!$('#tabs').tabs('exists',title)){
		$('#tabs').tabs('add',{
			title:title,
			content:'<iframe name="mainFrame" title="'+title+'" scrolling="auto" frameborder="0"  src="'+url+'" style="width:100%;height:99%;"></iframe>',
			closable:true,
			width:$('#mainPanle').width()-10,
			height:$('#mainPanle').height()-26
		});
	}else{
			$('#tabs').tabs('select',title);
	}
}