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
	export class Proxy {
	    url: string;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new Proxy(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.type = source["type"];
	    }
	}
	export class Queue {
	    username: string;
	    startTime: number;
	    endTime: number;
	    infinite: boolean;
	    useProxies: boolean;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new Queue(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.username = source["username"];
	        this.startTime = source["startTime"];
	        this.endTime = source["endTime"];
	        this.infinite = source["infinite"];
	        this.useProxies = source["useProxies"];
	        this.status = source["status"];
	    }
	}

}

