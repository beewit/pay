var laytpl;
var backUrl, token, fid;
var mtTpl = funcList.innerHTML, funcView = document.getElementById('funcView');
var fcTpl = funcCharge.innerHTML, funcChrageView = document.getElementById('funcChrageView');


$(function () {
    backUrl = LinkUrl.Request('backUrl') || false;
    if (backUrl) {
        Cookies.set('backUrl', backUrl, {expires: 7})
    }
    token = LinkUrl.Request('token') || false;
    if (token) {
        Cookies.set('token', token, {expires: 7})
    } else {
        token = Cookies.get('token')
    }

    fid = LinkUrl.Request('fid') || "";
    if (fid == "") {
        if (backUrl) {
            parent.location.href = decodeURIComponent(backUrl)
        } else {
            parent.location.href = "http://www.tbqbz.com/"
        }
    }

    layui.use('laytpl', function () {
        laytpl = layui.laytpl;
        checkLogin()
    });

    // $(".type-vip").delegate(".vip-item", "click", function () {
    //     var t = $(this).attr("data-tag");
    //     $(".vip-item").removeClass("on");
    //     $(this).addClass("on");
    //     $(".type-price").hide();
    //     $("." + t + ".type-price").show();
    //     $("." + t + ".type-price").removeClass("on");
    //     $(".price-item").removeClass("on")
    //     $("." + t + ".type-price a[data-default=是]").addClass("on");
    //     $(".type-card").hide();
    //     $("." + t + ".type-card").show();
    //     $("." + t).show();
    //     createOrderComm()
    // });

    $(".type-price-box").delegate(".price-item", "click", function () {
        $(this).parent().find(".price-item").removeClass("on")
        $(this).addClass("on")
        createOrderComm()
    });

    $(".j-agreement").click(function () {
        $("#agreement").show();
    });

    $(".model-close").click(function () {
        $(".qtPayNew-model").hide();
    });

    $(".alipay_t").click(function () {
        var href = $(this).attr("data-href");
        if (href) {
            parent.location.href = href
        }
    });

});

function checkLogin() {
    loadFunc();
}

function loadFunc() {
    $.ajax({
        load_tip: false,
        url: "/member/type",
        data: {fid: fid},
        success: function (d) {
            if (d.ret == 200) {
                if (d.data.account != null && d.data.account != undefined) {
                    $("#username-logined").html(d.data.account.nickname);
                    $("#id").html(d.data.account.id);
                }
                laytpl(mtTpl).render(d.data.func, function (html) {
                    funcView.innerHTML = html;
                });
                var sumPrice = 0;
                $.each(d.data.func, function (i, item) {
                    sumPrice += parseFloat((item.price || 0))
                })
                d.data.sumPrice = sumPrice
                laytpl(fcTpl).render(d.data, function (html) {
                    funcChrageView.innerHTML = html;
                });
                createOrderComm()
            }
        }
    });
}


function createOrderComm() {
    if (token) {
        createOrder("支付宝")
        createOrder("微信")
    } else {
        layer.msg("请登陆后操作", {icon: 0})
    }
}

var ct;
var i = 0;


function createOrder(pt) {
    var funcIds = fid
    var $price = $(".type-price").find(".price-item.on");
    var priceId = $price.attr("data-id");
    var price = $price.attr("data-price");
    $(".order-price .o-p").html(price);
    $("#alipay-code").attr("src", "");
    $("#wechat-code").attr("src", "");
    $.ajax({
        load_tip: false,
        url: "/order/create/list",
        data: {funcId: funcIds, fcId: priceId, pt: pt},
        success: function (d) {
            if (d.ret == 200) {
                var tradeNo = d.data.tradeNo;
                if (pt == "支付宝") {
                    //启动定时查询任务
                    $("#alipay-code").attr("src", d.data.codeUrl);
                    $(".alipay_t").attr("data-href", d.data.getUrl);
                } else if (pt == "微信") {
                    $("#wechat-code").attr("src", d.data.codeUrl);
                }
                clearTimeout(ct)
                queryOrder(tradeNo);
            }
        }
    });
}


function queryOrder(tradeNo) {
    if (tradeNo) {
        $.ajax({
            load_tip: false,
            url: "/order/query",
            data: {tradeNo: tradeNo},
            success: function (d) {
                if (d.ret == 200) {
                    layer.msg("支付成功", {icon: 1}, function () {
                        var backUrl = Cookies.get('backUrl') || false;
                        if (backUrl) {
                            parent.location.href = decodeURIComponent(backUrl)
                        } else {
                            parent.location.href = "http://www.tbqbz.com/"
                        }
                    })
                }
                else if (d.ret == 404) {
                    timeOut(tradeNo)
                }
            }
        });
    }
}

function timeOut(tradeNo) {
    //2个小时后不再轮询
    if (i < 3600 * 2) {
        ct = setTimeout(function () {
            i++
            queryOrder(tradeNo)
            console.log(tradeNo)
        }, 1000)
    }
}