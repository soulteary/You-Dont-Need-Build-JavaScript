define(function (require) {
  var template = require("tpl!./container.html");

  var Layout = santd.Layout;
  var Menu = santd.Menu;
  var Icon = santd.Icon;
  var Breadcrumb = santd.Breadcrumb;

  return san.defineComponent({
    components: {
      "s-layout": Layout,
      "s-header": Layout.Header,
      "s-content": Layout.Content,
      "s-sider": Layout.Sider,
      "s-menu": Menu,
      "s-sub-menu": Menu.Sub,
      "s-menu-item": Menu.Item,
      "s-icon": Icon,
      "s-breadcrumb": Breadcrumb,
      "s-brcrumbitem": Breadcrumb.Item,
    },
    initData() {
      return {
        inlineCollapsed: false,
      };
    },
    toggleCollapsed() {
      this.data.set("inlineCollapsed", !this.data.get("inlineCollapsed"));
    },
    template: template,
  });
});
