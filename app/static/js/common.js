//全局版本变量
var version = "20170831163240";
var layer, player;

var win_name
layui.use('layer', function () {
    layer = layui.layer;
    player = layui.layer;// = window === window.top ? layer : (parent.window === parent.window.top) ? parent.layer :
    //(parent.parent.window === parent.parent.window.top ? parent.parent.layer : parent.parent.parent.layer);
    win_name = player.getFrameIndex(window.name);
    //页面中弹窗处理
    $("body").delegate("[data-layer]", "click", function () {
        var asyn = $(this).attr("data-layer-asyn") || false;
        if (!asyn) {
            openWin($(this), win_name);
        }
    });
});
$(function () {
    //if (location.hostname != "localhost" && location.hostname != "127.0.0.1") {
    // $(document).bind("contextmenu", function () {
    //     return false;
    // });
    // $(document).bind("selectstart", function () {
    //     return false;
    // });
    //$("body").attr("oncontextmenu", "return false").attr("onselectstart", "return false").attr("oncopy", "alert('不支持复制！');return false;");
    // }
});

function openWin($obj, win_name) {
    var handle = $obj.attr("data-layer-handle");
    var flog = true;
    if (handle) {
        flog = eval("flog=" + handle + "()");
    }
    if (!flog) {
        return;
    }
    var pkid = $("#pkid").val();
    var href = $obj.attr("data-layer");
    if (win_name != undefined) {
        if (href.indexOf("?") > -1) {
            href = href + "&winname=" + win_name;
        } else {
            href = href + "?winname=" + win_name;
        }
    }
    if (pkid != undefined && pkid != null && pkid != "" && href.indexOf("&id=") < 0 && href.indexOf("?id=") < 0) {
        if (href.indexOf("?") > -1) {
            href = href + "&id=" + pkid;
        } else {
            href = href + "?id=" + pkid;
        }
    }
    var title = $obj.attr("data-layer-title");
    var w = $obj.attr("data-layer-w");
    var h = $obj.attr("data-layer-h");
    w = w || "893px";
    h = h || "600px";
    var t = $obj.attr("data-layer-top");
    var l = $obj.attr("data-layer-left");
    t = t || "5%";
    l = l || "";
    var shade = $obj.attr("data-layer-shade");
    shade = shade != undefined ? [0.8, '#000000'] : false;
    var maxmin = $obj.attr("data-layer-maxmin");
    maxmin = maxmin == "false" ? false : true;
    var closeBtn = $obj.attr("data-layer-closeBtn");
    closeBtn = closeBtn == undefined ? 1 : parseInt(closeBtn);
    var security = $obj.attr("data-layer-iframe-security") != undefined;
    if (href != null && href != undefined && href != "") {
        player.open({
            type: 2,
            title: title,
            moveOut: true,
            resize: false,
            shade: shade,
            shadeClose: false,//如果shade存在，shadeClose控制点击弹层外区域关闭。
            maxmin: maxmin, //开启最大化最小化按钮
            area: [w, h],
            offset: [t, l],
            content: LinkUrl.ChangeUrlParas(href, "v", version),
            closeBtn: closeBtn,
            moveEnd: function (e) {
                var top = $(window.parent.document).find(e.prevObject.selector).offset().top;
                var left = $(window.parent.document).find(e.prevObject.selector).offset().left;
                var wh = $(window.parent.document).find(e.prevObject.selector).width();
                if (top < 0) {
                    $(window.parent.document).find(e.prevObject.selector).css("top", 0);
                } else if (top + 40 > parent.document.body.clientHeight) {
                    $(window.parent.document).find(e.prevObject.selector).css("top", (parent.document.body.clientHeight - 40) + "px");
                }
                if (left + wh < 40) {
                    $(window.parent.document).find(e.prevObject.selector).css("left", (40 - wh) + "px");
                }
                if (left + 40 > parent.document.body.clientWidth) {
                    $(window.parent.document).find(e.prevObject.selector).css("left", (parent.document.body.clientWidth - 40) + "px");
                }
            }, success: function (layero, index) {
                if (security) {
                    $(layero).find("iframe").attr("security", "restricted").attr("sandbox", "");
                }
            }
        });
    }
}

//关闭右键菜单
function closeRightMenu() {
    $("body").attr("oncontextmenu", "return false").attr("onselectstart", "return false");
    $("html").attr("oncontextmenu", "return false").attr("onselectstart", "return false");
}

//关闭弹出
function closeIframe() {
    var index = parent.layer.getFrameIndex(window.name);
    player.close(index);
}

//子弹窗调用父级弹窗的方法
var ExecuteParentMethod = function (methodName, param1, param2, param3) {
    try {
        switch (arguments.length) {
            case 1:
                eval("$(window.parent.document).find('#layui-layer-iframe' + (LinkUrl.Request('winname') || ''))[0].contentWindow." + methodName + "()");
                break;
            case 2:
                eval("$(window.parent.document).find('#layui-layer-iframe' + (LinkUrl.Request('winname') || ''))[0].contentWindow." + methodName + "('" + param1 + "')");
                break;
            case 3:
                eval("$(window.parent.document).find('#layui-layer-iframe' + (LinkUrl.Request('winname') || ''))[0].contentWindow." + methodName + "('" + param1 + "','" + param2 + "')");
                break;
            case 4:
                eval("$(window.parent.document).find('#layui-layer-iframe' + (LinkUrl.Request('winname') || ''))[0].contentWindow." + methodName + "('" + param1 + "','" + param2 + "','" + param3 + "')");
                break;
            default:

        }
    } catch (e) {
    }
};

//子弹窗调用Iframe的方法
var ExecuteIframeMethod = function (methodName, param1, param2, param3) {
    try {
        switch (arguments.length) {
            case 1:
                eval("$(window.parent.document).find('#hive-iframe')[0].contentWindow." + methodName + "()");
                break;
            case 2:
                eval("$(window.parent.document).find('#hive-iframe')[0].contentWindow." + methodName + "('" + param1 + "')");
                break;
            case 3:
                eval("$(window.parent.document).find('#hive-iframe')[0].contentWindow." + methodName + "('" + param1 + "','" + param2 + "')");
                break;
            case 4:
                eval("$(window.parent.document).find('#hive-iframe')[0].contentWindow." + methodName + "('" + param1 + "','" + param2 + "','" + param3 + "')");
                break;
            default:

        }
    } catch (e) {
    }
};


function LoadScript(url, callback) {
    var script = document.createElement("script");
    script.type = "text/javascript";
    if (typeof(callback) != "undefined") {
        if (script.readyState) {
            script.onreadystatechange = function () {
                if (script.readyState == "loaded" || script.readyState == "complete") {
                    script.onreadystatechange = null;
                    callback();
                }
            };
        } else {
            script.onload = function () {
                callback();
            };
        }
    }
    script.src = url;
    document.body.appendChild(script);
}

var LinkUrl = {
    Request: function (name) {
        var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)", "i");
        var search = window.location.search.substr(1);
        search = decodeURI(search);
        var r = search.match(reg);
        if (r != null) {
            return (r[2]);
        }
        return "";
    },
    ChangeUrlParas: function (url, ref, value) {
        var str = "";
        url = this.DeleteParas(url, "page");
        if (url.indexOf('?') != -1)
            str = url.substr(url.indexOf('?') + 1);
        else
            return url + "?" + ref + "=" + value;
        var returnurl = "";
        var setparam = "";
        var arr;
        var modify = "0";
        if (str.indexOf('&') != -1) {
            arr = str.split('&');
            for (var i = 0; i < arr.length; i++) {
                if (arr[i].split('=')[0] == ref) {
                    setparam = value;
                    modify = "1";
                }
                else {
                    setparam = arr[i].split('=')[1];
                }
                returnurl = returnurl + arr[i].split('=')[0] + "=" + setparam + "&";
            }
            returnurl = returnurl.substr(0, returnurl.length - 1);
            if (modify == "0")
                if (returnurl == str)
                    returnurl = returnurl + "&" + ref + "=" + value;
        }
        else {
            if (str.indexOf('=') != -1) {
                arr = str.split('=');
                if (arr[0] == ref) {
                    setparam = value;
                    modify = "1";
                }
                else {
                    setparam = arr[1];
                }
                returnurl = arr[0] + "=" + setparam;
                if (modify == "0")
                    if (returnurl == str)
                        returnurl = returnurl + "&" + ref + "=" + value;
            }
            else
                returnurl = ref + "=" + value;
        }
        return url.substr(0, url.indexOf('?')) + "?" + returnurl;
    },
    DeleteParas: function (url, ref) {
        var str = "";
        if (url.indexOf('?') != -1) {
            str = url.substr(url.indexOf('?') + 1);
        }
        else {
            return url;
        }
        var arr = "";
        var returnurl = "";
        if (str.indexOf('&') != -1) {
            arr = str.split('&');
            for (var i = 0; i < arr.length; i++) {
                '';
                if (arr[i].split('=')[0] != ref) {
                    returnurl = returnurl + arr[i].split('=')[0] + "=" + arr[i].split('=')[1] + "&";
                }
            }
            return url.substr(0, url.indexOf('?')) + "?" + returnurl.substr(0, returnurl.length - 1);
        }
        else {
            arr = str.split('=');
            if (arr[0] == ref) {
                return url.substr(0, url.indexOf('?'));
            }
            else {
                return url;
            }
        }
    }
};
var Request = GetRequest();

function GetRequest() {
    var url = location.search; //获取url中"?"符后的字串
    var theRequest = new Object();
    if (url.indexOf("?") != -1) {
        var str = url.substr(1);
        var strs = str.split("&");
        for (var i = 0; i < strs.length; i++) {
            theRequest[strs[i].split("=")[0]] = unescape(strs[i].split("=")[1]);
        }
    }
    return theRequest;
}


function encrypt(str, pwd) {
    if (pwd == null || pwd.length <= 0) {
        alert("Please enter a password with which to encrypt the message.");
        return null;
    }
    var prand = "";
    for (var i = 0; i < pwd.length; i++) {
        prand += pwd.charCodeAt(i).toString();
    }
    var sPos = Math.floor(prand.length / 5);
    var mult = parseInt(prand.charAt(sPos) + prand.charAt(sPos * 2) + prand.charAt(sPos * 3) + prand.charAt(sPos * 4) + prand.charAt(sPos * 5));
    var incr = Math.ceil(pwd.length / 2);
    var modu = Math.pow(2, 31) - 1;
    if (mult < 2) {
        alert("Algorithm cannot find a suitable hash. Please choose a different password. \nPossible considerations are to choose a more complex or longer password.");
        return null;
    }
    var salt = Math.round(Math.random() * 1000000000) % 100000000;
    prand += salt;
    while (prand.length > 10) {
        prand = (parseInt(prand.substring(0, 10)) + parseInt(prand.substring(10, prand.length))).toString();
    }
    prand = (mult * prand + incr) % modu;
    var enc_chr = "";
    var enc_str = "";
    for (var i = 0; i < str.length; i++) {
        enc_chr = parseInt(str.charCodeAt(i) ^ Math.floor((prand / modu) * 255));
        if (enc_chr < 16) {
            enc_str += "0" + enc_chr.toString(16);
        } else enc_str += enc_chr.toString(16);
        prand = (mult * prand + incr) % modu;
    }
    salt = salt.toString(16);
    while (salt.length < 8) salt = "0" + salt;
    enc_str += salt;
    return enc_str;
}

function decrypt(str, pwd) {
    if (str == null || str.length < 8) {
        alert("A salt value could not be extracted from the encrypted message because it's length is too short. The message cannot be decrypted.");
        return;
    }
    if (pwd == null || pwd.length <= 0) {
        alert("Please enter a password with which to decrypt the message.");
        return;
    }
    var prand = "";
    for (var i = 0; i < pwd.length; i++) {
        prand += pwd.charCodeAt(i).toString();
    }
    var sPos = Math.floor(prand.length / 5);
    var mult = parseInt(prand.charAt(sPos) + prand.charAt(sPos * 2) + prand.charAt(sPos * 3) + prand.charAt(sPos * 4) + prand.charAt(sPos * 5));
    var incr = Math.round(pwd.length / 2);
    var modu = Math.pow(2, 31) - 1;
    var salt = parseInt(str.substring(str.length - 8, str.length), 16);
    str = str.substring(0, str.length - 8);
    prand += salt;
    while (prand.length > 10) {
        prand = (parseInt(prand.substring(0, 10)) + parseInt(prand.substring(10, prand.length))).toString();
    }
    prand = (mult * prand + incr) % modu;
    var enc_chr = "";
    var enc_str = "";
    for (var i = 0; i < str.length; i += 2) {
        enc_chr = parseInt(parseInt(str.substring(i, i + 2), 16) ^ Math.floor((prand / modu) * 255));
        enc_str += String.fromCharCode(enc_chr);
        prand = (mult * prand + incr) % modu;
    }
    return enc_str;
}

function removeDecimalPoint(dp) {
    if (dp) {
        return dp.replace(".00", "")
    }
    return dp
}