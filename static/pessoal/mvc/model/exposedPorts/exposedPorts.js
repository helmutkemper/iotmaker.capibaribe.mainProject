var containerHostExposedPortsAddNewPortDialogWindowRef;
var containerConfigurationExposedPortsNumberRef;
var containerHostExposedPortsRef;

var containerHostExposedPortsItemsCount = -1;
var containerHostExposedPortsItemsIdToAdd = -1;

golbalOnInit.push(function(){
  $("#containerHostExposedPorts").kendoMultiSelect({
    filter: "contains",
    placeholder: "",
    itemTemplate: kendo.template($("#containerHostExposedPortsTemplate").html()),
    noDataTemplate: kendo.template($("#containerHostExposedPortsNoDataTemplate").html()),
    footerTemplate: kendo.template($("#containerHostExposedPortsFooterTemplate").html()),
    dataTextField: "Value",
    dataValueField: "Id",
    dataSource: new kendo.data.DataSource({
      schema: {
        model: {
          id: "Id",
          fields: {
            Id: { type: "number" },
            Value: { type: "string" },
            ImageName: { type: "string" }
          }
        }
      }
    })
  });
  containerHostExposedPortsRef = $("#containerHostExposedPorts").data("kendoMultiSelect");
  containerHostExposedPortsAddNewPortDialogWindowRef = $("#containerHostExposedPortsAddNewPort");
});

function containerHostExposedPortsAddNewPort(){
  containerHostExposedPortsAddNewPortDialogWindowRef.kendoDialog({
    title: "Expose port from container",
    content: kendo.template($("#containerCreateTemplateExposedPortsAddNewPort").html()),
    visible: false,
    modal: true,
    close: function(){

    },
    open: function(){
      $("#ExposedPortsNumber").kendoNumericTextBox({ format: "#" });
      containerConfigurationExposedPortsNumberRef = $("#ExposedPortsNumber").data("kendoNumericTextBox");

      $("#ExposedPortsProtocol").kendoDropDownList();

      setTimeout( function(){ containerConfigurationExposedPortsNumberRef.focus(); }, 1000)
    },
    actions: [
      {
        text: "Close"
      },
      {
        text: "Add",
        action: function(e){
          containerHostExposedPortsAddNewPortFunction();
          return false;
        }
      },
      {
        text: "Add and close",
        action: function(e){
          containerHostExposedPortsAddNewPortFunction();
        },
        primary: true
      }
    ]
  });
  containerHostExposedPortsAddNewPortDialogWindowRef.data("kendoDialog").open();
}

function containerHostExposedPortsAddNewPortFunction(){
  let imageName = containerConfigurationImageNameRef.text();

  // Procura por um nome de container
  if( imageName == "" ) {
    $("#dialog").kendoDialog({
      modal: true,
      visible: true,
      title: "Configuration error",
      content: "Please, select an image name first.",
      width: 400,
      actions: [
        { text: "OK", primary: true, action: function(){ containerConfigurationImageNameRef.open(); } }
      ]
    });
    return;
  }

  // Procura por uma porta válida
  if( containerConfigurationExposedPortsNumberRef.value() == null ){
    $("#dialog").kendoDialog({
      modal: true,
      visible: true,
      title: "Configuration error",
      content: "Please, select a valid port number.",
      width: 400,
      actions: [
        { text: "OK", primary: true, action: function(){ setTimeout( function(){ containerConfigurationExposedPortsNumberRef.focus(); }, 1000); } }
      ]
    });
    return;
  }

  let dataSource = containerHostExposedPortsRef.dataSource;
  containerHostExposedPortsItemsCount += 1;
  containerHostExposedPortsItemsIdToAdd = containerHostExposedPortsItemsCount;

  let dataInDataSource = dataSource.data();
  let pass = true;

  for (var i in dataInDataSource) {
    if( isNaN( i ) ){
      break
    }
    i = parseInt(i);
    if (dataInDataSource[i].Value === $("#ExposedPortsNumber").val() + "/" + $("#ExposedPortsProtocol").val()) {
      containerHostExposedPortsItemsIdToAdd = dataInDataSource[i].Id;
      pass = false;
      break;
    }
  }

  if( pass === true ) {
    dataSource.add({
      Id: containerHostExposedPortsItemsCount,
      Value: $("#ExposedPortsNumber").val() + "/" + $("#ExposedPortsProtocol").val(),
      ImageName: containerConfigurationImageNameRef.text()
    });
  }

  dataSource.one("requestEnd", function(args) {
    if (args.type !== "create") {
      return;
    }

    dataSource.one("sync", function() {
      containerHostExposedPortsRef.value(containerHostExposedPortsRef.value().concat([containerHostExposedPortsItemsIdToAdd]));
    });

    $("#ExposedPortsNumber").data("kendoNumericTextBox").value("");
  });

  dataSource.sync();
}