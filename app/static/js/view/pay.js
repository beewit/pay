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
        createAlipayOrder()
    });

    $(".type-price-box").delegate(".price-item", "click", function () {
        $(this).parent().find(".price-item").removeClass("on")
        $(this).addClass("on")
        createAlipayOrder()
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

    queryOrder()
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
                createAlipayOrder()
            }
        }
    });
}

var tradeNo;
var i = 0;

function createAlipayOrder() {
    var $vip = $(".vip-item.on");
    var vipId = $vip.attr("data-id")
    var $price = $(".type-price." + $vip.attr("data-tag")).find(".price-item.on");
    var priceId = $price.attr("data-id");
    var price = $price.attr("data-price");
    $(".order-price .o-p").html(price);
    $.ajax({
        load_tip: false,
        url: "/order/create",
        data: {mtId: vipId, mtcId: priceId, pt: "支付宝"},
        success: function (d) {
            if (d.ret == 200) {
                tradeNo = d.data.tradeNo
                //启动定时查询任务
                $("#alipay-code").attr("src", d.data.codeUrl)
                $(".alipay_t").attr("data-href", d.data.getUrl)
            }
        }
    });
}


function queryOrder() {
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
                    timeOut()
                }
            }
        });
    } else {
        timeOut()
    }
}

function timeOut() {
    if (i < 3600 * 2) {
        setTimeout(function () {
            i++
            queryOrder()
        }, 1000)
    } else {
        //2个小时后不再轮询
    }
}