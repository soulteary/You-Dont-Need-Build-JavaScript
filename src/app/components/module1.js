define(function (require) {
  var template = require("tpl!./module1.html");

  var Message = santd.Message;
  var DatePicker = santd.DatePicker;

  return san.defineComponent({
    initData() {
      return {
        date: null,
      };
    },
    components: {
      "s-datepicker": DatePicker,
      "router-link": sanRouter.Link,
    },
    handleChange({ date }) {
      Message.info(`您选择的日期是: ${date ? date.format("YYYY-MM-DD") : "未选择"}`);
      this.data.set("date", date);
    },
    getDate(date) {
      return date ? date.format("YYYY-MM-DD") : "未选择";
    },
    template: template,
  });
});
