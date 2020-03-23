package main

import (
	"net/http"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "fmt"
    "io/ioutil"
    "context"
    "log"
    "encoding/json"
    "crypto/rsa"
    "crypto/rand"
)

type UserEntry struct {
	ID			*primitive.ObjectID	`bson:"_id,omitempty"`
	Name	 	string				`bson:"name"`
	Password 	string				`bson:"password"`
	Private_key *rsa.PrivateKey		`bson:"private_key"`
	Public_key	*rsa.PublicKey		`bson:"public_key"`
	Group 		[]string			`bson:"group",omitempty`
}

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
		response, err := encryptTheMessage(decodedRequest)
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
		response, err := decryptTheMessage(decodedRequest)
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
		response, err := addToGroup(decodedRequest)
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
		response, err := removeFromGroup(decodedRequest)
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
		response, err := logIn(decodedRequest)
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
		success, err := registerUser(decodedRequest)
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
	
	var result UserEntry
	filter := bson.M{"name": username, "password": password}
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Printf("User does not exist or there was another problem: %v", err)
		return Response{}, err
	}
	recipients := result.Group
	var ciphertext string
	for _, recipient := range recipients {
		encryptedMessage, err := encryptSingle(recipient, req.Message)
		if err != nil {
			log.Printf("Problem encrypting the message: %v", err)
			return Response{}, err
		}
		ciphertext = ciphertext+ " "+recipient+":"+encryptedMessage	
	}
	response := Response{
		Message: ciphertext,
	}
	return response, nil
}

//Encrypts the message for a single user
func encryptSingle(username, plaintext string) (string, error) {	
	var result UserEntry
	filter := bson.M{"name": username}
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Printf("User does not exist or there was another problem: %v", err)
		return "", err
	}
	
	publicKey := result.Public_key
	ciphertextBytes, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(plaintext))
	if err != nil {
		return "", err
	}
	ciphertext := string(ciphertextBytes)
	return ciphertext, nil

}

//Decrypts the message and sends back the plaintext
func decryptTheMessage(req Request) (Response, error) {
	username := req.User
	password := req.Password
	
	var result UserEntry
	filter := bson.M{"name": username, "password": password}
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Printf("User does not exist or there was another problem: %v", err)
		return Response{}, err
	}
	
	privateKey := result.Private_key
	plaintextBytes, err := privateKey.Decrypt(rand.Reader, []byte(req.Message), nil)
	if err != nil {
		log.Printf("Problem decrypting the message: %v", err)
	}
	plaintext := string(plaintextBytes)
	response := Response{
		Message: plaintext,
	}
	return response, nil
	
}

//Adds a user to the group and updates BD before returning new list to extension
func addToGroup(req Request) (GroupListResponse, error) {
	username := req.User
	password := req.Password
	
	var result UserEntry
	filter := bson.M{"name": username, "password": password}
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Printf("User does not exist or there was another problem: %v", err)
		return GroupListResponse{}, err
	}

	recipients := append(result.Group, req.Message)
	
	replacement := bson.M{
		"name" : username,
		"password" : password,
		"private_key" : result.Private_key,
		"public_key" : result.Public_key,
		"group" : recipients}
	var replacedDoc bson.M 
	err = collection.FindOneAndReplace(context.Background(), filter, replacement).Decode(&replacedDoc)
	if err != nil {
	    log.Printf("Problem replacing the document: %v", err)
	    return GroupListResponse{}, err
	}
	response := GroupListResponse{
		GroupMembers : recipients,
	}
	return response, nil
}

//Removes a user from the group and updates BD before returning new list to extension
func removeFromGroup(req Request) (GroupListResponse, error) {
	username := req.User
	password := req.Password
	
	var result UserEntry
	filter := bson.M{"name": username, "password": password}
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Printf("User does not exist or there was another problem: %v", err)
		return GroupListResponse{}, err
	}

	toRemove := req.Message
	var newRecipients []string
	for _, recipient := range result.Group {
		if recipient == toRemove {
			continue
		} else {
			newRecipients = append(newRecipients, recipient)
		}
	}
	
	replacement := bson.M{
		"name" : username,
		"password" : password,
		"private_key" : result.Private_key,
		"public_key" : result.Public_key,
		"group" : newRecipients}
	var replacedDoc bson.M 
	err = collection.FindOneAndReplace(context.Background(), filter, replacement).Decode(&replacedDoc)
	if err != nil {
	    log.Printf("Problem replacing the document: %v", err)
	    return GroupListResponse{}, err
	}
	response := GroupListResponse{
		GroupMembers : newRecipients,
	}
	return response, nil
}

//Check if user exists and has exact password, and then return his group list
func logIn(req Request) (GroupListResponse, error) {
	username := req.User
	password := req.Password
	
	var result UserEntry
	filter := bson.M{"name": username, "password": password}
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Printf("User does not exist or there was another problem: %v", err)
		return GroupListResponse{}, err
	}
	response := GroupListResponse{
		GroupMembers : result.Group,
	}
	return response, nil
}

//Register the user with given username and password to DB
func registerUser(req Request) (bool, error) {
	username := req.User
	password := req.Password
	
	var result UserEntry
	filter := bson.D{{"name",username}}
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err!= nil {
		if err == mongo.ErrNoDocuments {
			privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				log.Printf("Private key generation failed: %v", err)
			}
			publicKey := privateKey.Public()
			_, err = collection.InsertOne(context.Background(), bson.M{"name": username, "password": password, "private_key": privateKey, "public_key": publicKey})
			if err != nil {
				log.Printf("Problem registering the user: %v", err)
				return false, err
			}
			return true, nil
		} else {
			log.Printf("Problem registering the user: %v", err)
			return false, err
		}
	}
	log.Printf("User already exists")
	return false, nil
}