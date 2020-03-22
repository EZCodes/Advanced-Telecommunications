package main

import (
	"net/http"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
    "fmt"
    "io/ioutil"
    "context"
    "log"
    "encoding/json"
    "crypto/rsa"
    "crypto"
    "crypto/rand"
)

type Request struct {
	Type 		string		`json:type`
	User 		string		`json:user`
	Password 	string		`json:password`
	Message 	string		`json:message,omitempty`
	Recipients 	[]string	`json:recipients,omitempty`	
}

type Response struct {
	Message 	string		`json:message`
}

type GroupListResponse struct {
	GroupMembers []string	`json:members,omitempty`
}

var collection *(mongo.Collection)

func main(){
	// get mongoDB username and password
	m_username, err := ioutil.ReadFile("src/backend/username.txt") // file with just mongoDB username in it
	if err != nil {
    	log.Fatal(err)
    }
	m_password, err := ioutil.ReadFile("src/backend/password.txt") // file with just mongoDB password in it
	if err != nil {
    	log.Fatal(err) 
    }
	URI := "mongodb+srv://" + string(m_username) + ":" + string(m_password) + "@telecomms-mkx7q.mongodb.net/test?retryWrites=true&w=majority"
	
	// Set MongoDB client options
	clientOptions := options.Client().ApplyURI(URI)
	mongo_client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
	    log.Fatal(err)
	}
	// Check the connection
	err = mongo_client.Ping(context.Background(), nil)
	if err != nil {
	    log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	
	collection = mongo_client.Database("Telecomms").Collection("Userbase")
	
	log.Fatal(http.ListenAndServe(":420", http.HandlerFunc(requestHandler)))
	
}

//Handles the requests from the extension
func requestHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("Request received!")
	var decodedRequest Request
	
	err := json.NewDecoder(req.Body).Decode(&decodedRequest)
	if err != nil {
		log.Printf("JSON decoding failed!")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if decodedRequest.Type == "encrypt" {
		response, err := encryptTheMessage(request)
		if err != nil {
			log.Printf("Encrypting failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Printf("Marshalling failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	} else if decodedRequest.Type == "decrypt" {
		response, err := decryptTheMessage(request)
		if err != nil {
			log.Printf("Decrypting failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Printf("Marshalling failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	} else if decodedRequest.Type == "add" {
		response, err = addToGroup(request)
		if err != nil {
			log.Printf("Adding failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Printf("Marshalling failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	} else if decodedRequest.Type == "remove" {
		response, err = removeFromGroup(request)
		if err != nil {
			log.Printf("Removing failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Printf("Marshalling failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	} else if decodedRequest.Type == "login" {
		response, err := logIn(request)
		if err != nil {
			log.Printf("Login failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Printf("Marshalling failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	} else if decodedRequest.Type == "register" {
		success, err := registerUser(request)
		if err != nil {
			log.Printf("Registration failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if success {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	
	}	
}

//Encrypts the message and sends back ciphertext/s
func encryptTheMessage(req Request) (Response, error) {
	username := req.User
	password := req.Password
	
	var result bson.M
	filter := bson.M{"name": username, "password": password}
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Printf("User does not exist or there was another problem: %v", err)
		return nil, err
	}
	recipients := result.group
	var ciphertext string
	for _, recipient := range recipients {
		encryptedMessage, err := encryptSingle(recipient, req.Message)
		if err != nil {
			log.Printf("Problem encrypting the message: %v", err)
			return nil, err
		}
		ciphertext = ciphertext+ " "+recipient+":"+encryptedMessage	
	}
	response := Response{
		Message: ciphertext
	}
	return response, nil
}

//Encrypts the message for a single user
func encryptSingle(user, plaintext string) string, error {
	username := req.User
	password := req.Password
	
	var result bson.M
	filter := bson.M{"name": username, "password": password}
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Printf("User does not exist or there was another problem: %v", err)
		return nil, err
	}
	
	publicKey := result.public_key
	ciphertextBytes, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, byte(plaintext))
	if err != nil {
		return err
	}
	ciphertext := string(ciphertextBytes)
	return ciphertext

}

//Decrypts the message and sends back the plaintext
func decryptTheMessage(req Request) (Response, error) {
	username := req.User
	password := req.Password
	
	var result bson.M
	filter := bson.M{"name": username, "password": password}
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Printf("User does not exist or there was another problem: %v", err)
		return nil, err
	}
	
	privateKey := result.private_key
	plaintextBytes, err := privateKey.Decrypt(rand.Reader, byte(req.Message))
	if err != nil {
		log.Printf("Problem decrypting the message: %v", err)
	}
	plaintext := string(plaintextBytes)
	response := Response{
		Message: plaintext
	}
	return response, nil
	
}

//Adds a user to the group and updates BD before returning new list to extension
func addToGroup(req Request) (GroupListResponse, error) {
	username := req.User
	password := req.Password
	
	var result bson.M
	filter := bson.M{"name": username, "password": password}
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Printf("User does not exist or there was another problem: %v", err)
		return nil, err
	}

	recipients := append(result.group, req.Message)
	
	replacement := bson.M{
		"name" : username,
		"password" : password,
		"private_key" : result.private_key,
		"public_key" : result.public_key,
		"group" : recipients}
	var bson.M replacedDoc
	err := coll.FindOneAndReplace(context.Background(), filter, replacement).Decode(&replacedDoc)
	if err != nil {
	    log.Printf("Problem replacing the document: %v", err)
	    return nil, err
	}
	response := GroupListResponse{
		GroupMembers = recipients
	}
	return response, nil
}

//Removes a user from the group and updates BD before returning new list to extension
func removeFromGroup(req Request) (GroupListResponse, error) {
	username := req.User
	password := req.Password
	
	var result bson.M
	filter := bson.M{"name": username, "password": password}
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Printf("User does not exist or there was another problem: %v", err)
		return nil, err
	}

	recipients := result.group

	toRemove := req.Message
	var newRecipients []string
	for _, recipient := range recipients.GroupMembers {
		if recipient == toRemove {
			continue
		} else {
			newRecipients = append(newRecipients, recipient)
		}
	}
	
	replacement := bson.M{
		"name" : username,
		"password" : password,
		"private_key" : result.private_key,
		"public_key" : result.public_key,
		"group" : newRecipients}
	var bson.M replacedDoc
	err := coll.FindOneAndReplace(context.Background(), filter, replacement).Decode(&replacedDoc)
	if err != nil {
	    log.Printf("Problem replacing the document: %v", err)
	    return nil, err
	}
	response := GroupListResponse{
		GroupMembers = newRecipients
	}
	return response, nil
}

//Check if user exists and has exact password, and then return his group list
func logIn(req Request) (GroupListResponse, error) {
	username := req.User
	password := req.Password
	
	var result bson.M
	filter := bson.M{"name": username, "password": password}
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Printf("User does not exist or there was another problem: %v", err)
		return nil, err
	}
	response := GroupListResponse{
		GroupMembers = result.group
	}
	return response, nil
}

//Register the user with given username and password to DB
func registerUser(req Request) (bool, error) {
	username := req.User
	password := req.Password
	
	var result bson.M
	filter := bson.D{{"name",username}}
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err!= nil {
		if err == mongo.ErrNoDocuments {
			privateKey := rsa.GenerateKey(rand.Reader, 2048)
			publicKey := privateKey.Public()
			_, err := collection.InsertOne(context.Background(), bson.M{"name": username, "password": password, "private_key": privateKey, "public_key": publicKey})
			if err != nil {
				log.Printf("Problem registering the user: %v", err)
				return false, err
			}
			return true, nil
		} else {
			log.Printf("Problem registering the user: %v", err)
			return false, err
		}
		log.Printf("User already exists")
		return false, nil
	}
}