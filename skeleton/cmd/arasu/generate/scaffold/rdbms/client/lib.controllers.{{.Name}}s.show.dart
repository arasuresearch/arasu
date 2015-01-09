part of {{.App.Name}}.controllers.{{.Name}}s;
@Component(selector: '{{.Name}}s-show',  
  templateUrl: 'packages/{{.App.Name}}/views/{{.Name}}s/show.html' )
class {{.Cname}}sShow extends Cntr{
  {{.Cname}} record;

  {{.Cname}}sShow(RouteProvider routeProvider,Db db,NgRoutingHelper locationService) :
    super({{.Cname}},routeProvider,db,locationService) {
    ready.then((_){
      record = store.find(int.parse(Id, onError: (_) => 0));
    }).catchError((e)=>Logger.root.shout);
  }
  void destroy(var record){
    store.remove(record);
    locationService.router.gotoUrl('/{{.Name}}s');
  }
}
