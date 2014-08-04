part of components.{{.Name}}s;
@Component(selector: 'new',  templateUrl: 'packages/{{.App.Name}}/components/{{.Name}}s/form.html', publishAs: 'C')
class New extends Cntr {
  {{.Cname}} record;
  bool formSubmitted;

  New(RouteProvider routeProvider,Db db,NgRoutingHelper locationService) :
    super({{.Cname}},routeProvider,db,locationService) {
    ready.then((_){
      record = store.New();
      record.updateClone();
    }).catchError((e)=>Logger.root.shout(e));
  }

  void submit(form) {
    formSubmitted = true;
    if(form.invalid) return;
    if (record.hasChanged) {
      store.add(record).then((_){
        locationService.router.gotoUrl('/');
      }).catchError((e)=>Logger.root.shout(e));
    }
  }
}