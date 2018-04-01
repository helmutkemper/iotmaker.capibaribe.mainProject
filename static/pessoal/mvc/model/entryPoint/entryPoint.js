golbalOnInit.push(function(){
  containerHostEntryPointAddNewEntryPointDialogWindowRef = $("#containerHostEntryPointAddNewEntryPoint");

  $("#containerHostEntryPoint").kendoMultiSelect({
    filter: "contains",
    placeholder: "",
    itemTemplate: kendo.template($("#containerHostEntryPointTemplate").html()),
    noDataTemplate: kendo.template($("#containerHostEntryPointNoDataTemplate").html()),
    footerTemplate: kendo.template($("#containerHostEntryPointFooterTemplate").html()),
    dataTextField: "EntryPoint",
    dataValueField: "Id",
    dataSource: new kendo.data.DataSource({
      schema: {
        model: {
          id: "Id",
          fields: {
            Id: { type: "number" },
            EntryPoint: { type: "string" }
          }
        }
      }
    })
  });
  containerHostEntryPointRef = $("#containerHostEntryPoint").data("kendoMultiSelect");
});

var containerConfigurationEntryPointNameRef;
var containerHostEntryPointRef;
var containerHostEntryPointAddNewEntryPointDialogWindowRef;

var containerHostEntryPointItemsCount = -1;
var containerHostEntryPointItemsIdToAdd = -1;

function containerHostEntryPointAddNewEntryPoint(){
  containerHostEntryPointAddNewEntryPointDialogWindowRef.kendoDialog({
    title: "Entry point commands from container",
    content: kendo.template($("#containerCreateTemplateEntryPointAddNewEntryPoint").html()),
    visible: false,
    modal: true,
    close: function(){

    },
    open: function(){
      containerConfigurationEntryPointNameRef = $("#EntryPointName");

      setTimeout( function(){ containerConfigurationEntryPointNameRef.focus(); }, 1000)
    },
    actions: [
      {
        text: "Close"
      },
      {
        text: "Add",
        action: function(e){
          containerHostEntryPointAddNewEntryPointFunction();
          return false;
        }
      },
      {
        text: "Add and close",
        action: function(e){
          containerHostEntryPointAddNewEntryPointFunction();
        },
        primary: true
      }
    ]
  });
  containerHostEntryPointAddNewEntryPointDialogWindowRef.data("kendoDialog").open();
}

function containerHostEntryPointAddNewEntryPointFunction(){


  let dataSource = containerHostEntryPointRef.dataSource;
  containerHostEntryPointItemsCount += 1;
  containerHostEntryPointItemsIdToAdd = containerHostEntryPointItemsCount;

  let dataInDataSource = dataSource.data();
  let pass = true;

  for (var i in dataInDataSource) {
    if( isNaN( i ) ){
      break
    }
    i = parseInt(i);
    if (dataInDataSource[i].EntryPoint === $("#EntryPointName").val()) {
      containerHostEntryPointItemsIdToAdd = dataInDataSource[i].Id;
      pass = false;
      break;
    }
  }

  if( pass === true ) {
    dataSource.add({
      Id: containerHostEntryPointItemsCount,
      EntryPoint: $("#EntryPointName").val()
    });
  }

  dataSource.one("requestEnd", function(args) {
    if (args.type !== "create") {
      return;
    }

    dataSource.one("sync", function() {
      containerHostEntryPointRef.value(containerHostEntryPointRef.value().concat([containerHostEntryPointItemsIdToAdd]));
    });

    $("#EntryPointName").val("");
  });

  dataSource.sync();
}
