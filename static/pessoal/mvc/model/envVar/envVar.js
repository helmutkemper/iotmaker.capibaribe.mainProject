golbalOnInit.push(function(){
  containerHostEnvAddNewPortDialogWindowRef = $("#containerHostEnvAddNewEnvVar");

  $("#containerHostEnv").kendoMultiSelect({
    filter: "contains",
    placeholder: "",
    itemTemplate: kendo.template($("#containerHostEnvTemplate").html()),
    noDataTemplate: kendo.template($("#containerHostEnvNoDataTemplate").html()),
    footerTemplate: kendo.template($("#containerHostEnvFooterTemplate").html()),
    dataTextField: "EnvVar",
    dataValueField: "Id",
    dataSource: new kendo.data.DataSource({
      schema: {
        model: {
          id: "Id",
          fields: {
            Id: { type: "number" },
            EnvVar: { type: "string" }
          }
        }
      }
    })
  });
  containerHostEnvRef = $("#containerHostEnv").data("kendoMultiSelect");
});

var containerConfigurationEnvVarNameRef;
var containerHostEnvRef;
var containerHostEnvAddNewPortDialogWindowRef;

var containerHostEnvItemsCount = -1;
var containerHostEnvItemsIdToAdd = -1;

function containerHostEnvVarAddNewEnvVar(){
  containerHostEnvAddNewPortDialogWindowRef.kendoDialog({
    title: "Environments vars from container",
    content: kendo.template($("#containerCreateTemplateEnvVarAddNewEnvVar").html()),
    visible: false,
    modal: true,
    close: function(){

    },
    open: function(){
      containerConfigurationEnvVarNameRef = $("#EnvVarName");

      setTimeout( function(){ containerConfigurationEnvVarNameRef.focus(); }, 1000)
    },
    actions: [
      {
        text: "Close"
      },
      {
        text: "Add",
        action: function(e){
          containerHostEnvAddNewEnvVarFunction();
          return false;
        }
      },
      {
        text: "Add and close",
        action: function(e){
          containerHostEnvAddNewEnvVarFunction();
        },
        primary: true
      }
    ]
  });
  containerHostEnvAddNewPortDialogWindowRef.data("kendoDialog").open();
}

function containerHostEnvAddNewEnvVarFunction(){


  let dataSource = containerHostEnvRef.dataSource;
  containerHostEnvItemsCount += 1;
  containerHostEnvItemsIdToAdd = containerHostEnvItemsCount;

  let dataInDataSource = dataSource.data();
  let pass = true;

  for (var i in dataInDataSource) {
    if( isNaN( i ) ){
      break
    }
    i = parseInt(i);
    if (dataInDataSource[i].EnvVar === $("#EnvVarName").val()) {
      containerHostEnvItemsIdToAdd = dataInDataSource[i].Id;
      pass = false;
      break;
    }
  }

  if( pass === true ) {
    dataSource.add({
      Id: containerHostEnvItemsCount,
      EnvVar: $("#EnvVarName").val()
    });
  }

  dataSource.one("requestEnd", function(args) {
    if (args.type !== "create") {
      return;
    }

    dataSource.one("sync", function() {
      containerHostEnvRef.value(containerHostEnvRef.value().concat([containerHostEnvItemsIdToAdd]));
    });

    $("#EnvVarName").val("");
  });

  dataSource.sync();
}
