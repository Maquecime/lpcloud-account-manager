# What is this project about ?

I made this project as a part of my technical degree, the goal was to implement REST web services working together while being hosted either on Google App Engine or Heroku

These services aim to represent a bank account manager, with whom the users interact to make loan request

The services review the demand and, following an algorithm, approve or refuse it

If the demand is approved, the account of the requested user is credited, else, it is not.

# How to get it work

## GAE Services : 

* AccManager
* AppManager
* Databases (Google Datastore)

## Heroku Services :

* LoanApproval
* CheckAccount


## ON Google App Engine

Install google SDK : 

<https://cloud.google.com/sdk/docs/>

Install gcloud CLI : 

<https://cloud.google.com/sdk/docs?hl=fr#install_the_latest_cloud_tools_version_cloudsdk_current_version>

Start by initialzing gcloud / logging in 

    $ gcloud init
    $ gcloud auth login

Install app-engine for go
    
    $ sudo apt-get install google-cloud-sdk-app-engine-go 
    
Create gcloud project    

    $ gcloud projects create <your-project-name> --set-as-default
    
Check project details 

    $ gcloud projects describe lpcloud-project-1 
    
For these steps, you'll need to activate facturation on your project 

    $ gcloud app create --project=lpcloud-project-1 
    $ gcloud app deploy

To navigate between multiple projects, you'll also need this 

    
    $ gcloud config set project <your-project-name>
    

## On Heroku

Start by installing heroku cli :

    $ sudo snap install heroku --classic 

Log-In : 

    $ heroku login
    
Create the project :

    $ heroku create

Push to master : 

    $ git push heroku master
    
___
    
##### How2 Commit your local changes ? 

---




* Build the project :

    
    $ go build -o bin/<your-build-file-name> -v .
    
* Test it locally :

    
    $ heroku local
    
* Clean the unused modules :


    $ go mod tidy
    $ go mod vendor
    
* Stage your changes and commit them


    $ git add -A .
    $ git commit -m ""
    
* Push to master 


    $ git push heroku master
    
* open the app 

    
    $ heroku open
     
