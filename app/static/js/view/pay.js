var laytpl;
var mtTpl = memberTypeList.innerHTML, memberTypeView = document.getElementById('memberTypeView');
var mtcTpl = memberTypeCharge.innerHTML, memberTypeChrageView = document.getElementById('memberTypeChrageView');
layui.use('laytpl', function () {
    laytpl = layui.laytpl;
    loadMemberType();
});
$(function () {
    var backUrl = LinkUrl.Request('backUrl') || false;
    if (backUrl) {
        Cookies.set('backUrl', backUrl, {expires: 7})
    }
    $(".type-vip").delegate(".vip-item", "click", function () {
        var t = $(this).attr("data-tag");
        $(".vip-item").removeClass("on");
        $(this).addClass("on");
        $(".type-price").hide();
        $("." + t + ".type-price").show();
        $("." + t + ".type-price").removeClass("on");
        $(".price-item").removeClass("on")
        $("." + t + ".type-price a[data-default=是]").addClass("on");
        $(".type-card").hide();
        $("." + t + ".type-card").show();
        $("." + t).show();
        createOrder("支付宝")
        createOrder("微信")
    });

    $(".type-price-box").delegate(".price-item", "click", function () {
        $(this).parent().find(".price-item").removeClass("on")
        $(this).addClass("on")
        createOrder("支付宝")
        createOrder("微信")
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
            location.href = href
        }
    });

});

function loadMemberType() {
    $.ajax({
        load_tip: false,
        url: "/member/type",
        success: function (d) {
            if (d.ret == 200) {
                laytpl(mtTpl).render(d.data.memberType, function (html) {
                    memberTypeView.innerHTML = html;
                });
                laytpl(mtcTpl).render(d.data, function (html) {
                    memberTypeChrageView.innerHTML = html;
                });
                createOrder("支付宝")
                createOrder("微信")
            }
        }
    });
}

var ct;
var i = 0;

function createOrder(pt) {
    var $vip = $(".vip-item.on");
    var vipId = $vip.attr("data-id")
    var $price = $(".type-price." + $vip.attr("data-tag")).find(".price-item.on");
    var priceId = $price.attr("data-id");
    var price = $price.attr("data-price");
    $(".order-price .o-p").html(price);
    $("#alipay-code").attr("src", "");
    $("#wechat-code").attr("src", "");
    $.ajax({
        load_tip: false,
        url: "/order/create",
        data: {mtId: vipId, mtcId: priceId, pt: pt},
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
                    if (d.data.type == "会员套餐") {
                        if (d.data.member_type_id == 1) {
                            location.href = "/app/page/notify/success.html"
                        } else {
                            location.href = "/app/page/notify/success-company.html"
                        }
                    } else {
                        layer.msg("支付成功", {icon: 1}, function () {
                            var backUrl = Cookies.get('backUrl') || false;
                            if (backUrl) {
                                location.href = backUrl
                            } else {
                                location.href = "http://www.tbqbz.com/"
                            }
                        })
                    }
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