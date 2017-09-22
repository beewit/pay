(function ($) {
    var player = layer;
    var _ajax = $.ajax;
    $.ajax = function (options) {
        var defaults = {
            load_tip: true,
            err_tip: true,
            suc_tip: true,
            type: "post",
            dataType: "json",
            url: null
        };
        var opt = $.extend(defaults, options);
        var fn = {
            error: function (XMLHttpRequest, textStatus, errorThrown) {},
            success: function (data, textStatus) {},
            beforeSend: function (XMLHttpRequest) {
                XMLHttpRequest.setRequestHeader('key', "valid-key");
                XMLHttpRequest.setRequestHeader('Authorization', "valid-key");
            },
            complete: function (XMLHttpRequest, textStatus) {

            }
        };
        if (opt.error) {
            fn.error = opt.error;
        }
        if (opt.success) {
            fn.success = opt.success;
        }
        var loadi;
        if (opt.load_tip) {
            loadi = player.msg('正在加载..', {
                icon: 6,
                time: -1
            })
        }
        var _opt = $.extend(opt, {
            beforeSend: function (XMLHttpRequest) {
                fn.beforeSend(XMLHttpRequest);
            },
            error: function (XMLHttpRequest, textStatus, errorThrown) {
                if (opt.err_tip) {
                    player.alert('网络异常刷新重试', {
                        title: "系统提示",
                        icon: 2,
                        skin: 'layer-ext-moon'
                    });
                }
                //错误方法增强处理
                fn.error(XMLHttpRequest, textStatus, errorThrown);
            },
            success: function (json, textStatus) {
                if (opt.err_tip) {
                    player.close(loadi);
                    if (json.ret != 200 && json.ret != 404 && json.ret != undefined && json.ret != null) {
                        if (json.ret == 401) {
                            player.msg(json.msg, {
                                icon: 0,
                                time: 1500
                            }, function () {
                                var backUrl = Cookies.get("backUrl");
                                if (backUrl) {
                                    parent.location.href = decodeURIComponent(backUrl);
                                }
                            });
                        } else {
                            //错误统一拦截提示
                            player.alert(json.msg, {
                                title: "系统提示",
                                icon: 8,
                                skin: 'layer-ext-moon'
                            });
                        }
                        return;
                    }
                    //成功回调方法增强处理
                    fn.success(json, textStatus);
                } else {
                    fn.success(json, textStatus);
                }
            },
            complete: function (XMLHttpRequest, textStatus) {
                fn.complete(XMLHttpRequest, textStatus);
            }
        });
        _ajax(_opt);
    };
})(jQuery);