golbalOnInit.push(function(){
  containerHostLabelsAddNewLabelsDialogWindowRef = $("#containerHostLabelsAddNewLabels");

  $("#containerHostLabels").kendoMultiSelect({
    filter: "contains",
    placeholder: "",
    itemTemplate: kendo.template($("#containerHostLabelsTemplate").html()),
    noDataTemplate: kendo.template($("#containerHostLabelsNoDataTemplate").html()),
    footerTemplate: kendo.template($("#containerHostLabelsFooterTemplate").html()),
    dataTextField: "Labels",
    dataValueField: "Id",
    dataSource: new kendo.data.DataSource({
      schema: {
        model: {
          id: "Id",
          fields: {
            Id: { type: "number" },
            Labels: { type: "string" }
          }
        }
      }
    })
  });
  containerHostLabelsRef = $("#containerHostLabels").data("kendoMultiSelect");
});

var containerConfigurationLabelsNameRef;
var containerHostLabelsRef;
var containerHostLabelsAddNewLabelsDialogWindowRef;

var containerHostLabelsItemsCount = -1;
var containerHostLabelsItemsIdToAdd = -1;

function containerHostLabelsAddNewLabels(){
  containerHostLabelsAddNewLabelsDialogWindowRef.kendoDialog({
    title: "Labels from container",
    content: kendo.template($("#containerCreateTemplateLabelsAddNewLabels").html()),
    visible: false,
    modal: true,
    close: function(){

    },
    open: function(){
      containerConfigurationLabelsNameRef = $("#LabelsKey");

      setTimeout( function(){ containerConfigurationLabelsNameRef.focus(); }, 1000)
    },
    actions: [
      {
        text: "Close"
      },
      {
        text: "Add",
        action: function(e){
          containerHostLabelsAddNewLabelsFunction();
          return false;
        }
      },
      {
        text: "Add and close",
        action: function(e){
          containerHostLabelsAddNewLabelsFunction();
        },
        primary: true
      }
    ]
  });
  containerHostLabelsAddNewLabelsDialogWindowRef.data("kendoDialog").open();
}

function containerHostLabelsAddNewLabelsFunction(){

  if( $("#LabelsKey").val() === "" ){
    $("#dialog").kendoDialog({
      modal: true,
      visible: true,
      title: "Configuration error",
      content: "Please, enter a valid lavel key.",
      width: 400,
      actions: [
        {
          text: "OK",
          primary: true,
          action: function(){
            setTimeout( function(){
                $("#LabelsKey").focus();
            },
            1000);
          }
        }
      ]
    });
    return;
  }

  if( $("#LabelsKey").val().search("=") != -1 ){
    $("#dialog").kendoDialog({
      modal: true,
      visible: true,
      title: "Configuration error",
      content: "Please, dont't use character '='.",
      width: 400,
      actions: [
        {
          text: "OK",
          primary: true,
          action: function(){
            setTimeout( function(){
                $("#LabelsKey").focus();
              },
              1000);
          }
        }
      ]
    });
    return;
  }

  if( $("#LabelsName").val() === "" ){
    $("#dialog").kendoDialog({
      modal: true,
      visible: true,
      title: "Configuration error",
      content: "Please, enter a valid label value.",
      width: 400,
      actions: [
        {
          text: "OK",
          primary: true,
          action: function(){
            setTimeout( function(){
                $("#LabelsName").focus();
              },
              1000);
          }
        }
      ]
    });
    return;
  }

  if( $("#LabelsName").val().search("=") != -1 ){
    $("#dialog").kendoDialog({
      modal: true,
      visible: true,
      title: "Configuration error",
      content: "Please, dont't use character '='.",
      width: 400,
      actions: [
        {
          text: "OK",
          primary: true,
          action: function(){
            setTimeout( function(){
                $("#LabelsName").focus();
              },
              1000);
          }
        }
      ]
    });
    return;
  }

  let dataSource = containerHostLabelsRef.dataSource;
  containerHostLabelsItemsCount += 1;
  containerHostLabelsItemsIdToAdd = containerHostLabelsItemsCount;

  let dataInDataSource = dataSource.data();
  let pass = true;

  for (var i in dataInDataSource) {
    if( isNaN( i ) ){
      break
    }
    i = parseInt(i);
    if (dataInDataSource[i].Labels === $("#LabelsKey").val() + "=" + $("#LabelsName").val()) {
      containerHostLabelsItemsIdToAdd = dataInDataSource[i].Id;
      pass = false;
      break;
    }
  }

  if( pass === true ) {
    dataSource.add({
      Id: containerHostLabelsItemsCount,
      Labels: $("#LabelsKey").val() + "=" + $("#LabelsName").val()
    });
  }

  dataSource.one("requestEnd", function(args) {
    if (args.type !== "create") {
      return;
    }

    dataSource.one("sync", function() {
      containerHostLabelsRef.value(containerHostLabelsRef.value().concat([containerHostLabelsItemsIdToAdd]));
    });

    $("#LabelsKey").val("");
  });

  dataSource.sync();
}
