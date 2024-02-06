export namespace backendmanager {
	
	export class Account {
	    email: string;
	    password: string;
	    username: string;
	    type: string;
	    bearer: string;
	
	    static createFrom(source: any = {}) {
	        return new Account(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.email = source["email"];
	        this.password = source["password"];
	        this.username = source["username"];
	        this.type = source["type"];
	        this.bearer = source["bearer"];
	    }
	}

}

