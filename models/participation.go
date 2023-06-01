package models

import (
	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	ParticipationCreate struct {
		EventID string `json:"event_id" bson:"event_id"`
		Comment string `json:"comment" bson:"comment"`
	}

	ParticipationUpdate struct {
		ID      string `json:"id" bson:"_id"`
		Status  string `json:"status" bson:"status"`
		Comment string `json:"comment" bson:"comment"`
		//Confirmer UserInternal `json:"confirmer" bson:"confirmer"`
	}
	ParticipationDatabase struct {
		ID      string `json:"id" bson:"_id"`
		UserID  string `json:"user_id" bson:"user_id"`
		EventID string `json:"event_id" bson:"event_id"`
		Comment string `json:"comment" bson:"comment"`
		Status  string `json:"status" bson:"status"`
		CrewID  string `json:"crew_id" bson:"crew_id"`
		//Confirmer UserInternal   `json:"confirmer" bson:"confirmer"`
		Modified vmod.Modified `json:"modified" bson:"modified"`
	}
	Participation struct {
		ID      string `json:"id" bson:"_id"`
		UserID  string `json:"user_id" bson:"user_id"`
		User    User   `json:"user" bson:"user"`
		EventID string `json:"event_id" bson:"event_id"`
		Comment string `json:"comment" bson:"comment"`
		Status  string `json:"status" bson:"status"`
		Event   Event  `json:"event" bson:"event"`
		CrewID  string `json:"crew_id" bson:"crew_id"`
		Crew    Crew   `json:"crew" bson:"crew"`
		//Confirmer UserInternal   `json:"confirmer" bson:"confirmer"`
		Modified vmod.Modified `json:"modified" bson:"modified"`
	}
	UserParticipation struct {
		ID      string      `json:"id" bson:"_id"`
		EventID string      `json:"event_id" bson:"event_id"`
		Comment string      `json:"comment" bson:"comment"`
		Status  string      `json:"status" bson:"status"`
		Event   EventPublic `json:"event" bson:"event"`
		CrewID  string      `json:"crew_id" bson:"crew_id"`
		Crew    CrewName    `json:"crew" bson:"crew"`
		//Confirmer UserInternal   `json:"confirmer" bson:"confirmer"`
		Modified vmod.Modified `json:"modified" bson:"modified"`
	}
	EventParticipation struct {
		ID      string    `json:"id" bson:"_id"`
		UserID  string    `json:"user_id" bson:"user_id"`
		User    User      `json:"user" bson:"user"`
		EventID string    `json:"event_id" bson:"event_id"`
		Comment string    `json:"comment" bson:"comment"`
		Status  string    `json:"status" bson:"status"`
		Event   ListEvent `json:"event" bson:"event"`
		CrewID  string    `json:"crew_id" bson:"crew_id"`
		Crew    Crew      `json:"crew" bson:"crew"`
		//Confirmer UserInternal   `json:"confirmer" bson:"confirmer"`
		Modified vmod.Modified `json:"modified" bson:"modified"`
	}
	ParticipationMinimal struct {
		ID      string `json:"id" bson:"_id"`
		EventID string `json:"event_id" bson:"event_id"`
		Status  string `json:"status" bson:"status"`
	}

	ParticipationParam struct {
		ID string `param:"id"`
	}

	ParticipationQuery struct {
		ID       []string `query:"id" qs:"id"`
		EventID  []string `query:"event_id" qs:"event_id"`
		Comment  []string `query:"comment" bson:"comment"`
		Status   []string `query:"status" bson:"status"`
		UserId   []string `query:"user_id" bson:"user_id"`
		CrewName []string `query:"crew_name" bson:"crew_name"`
		CrewId   []string `query:"crew_id" qs:"crew_id"`
	}
	ParticipationStateRequest struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	ParticipationImport struct {
		DropsID string `json:"drops_id"`
		Comment string `json:"comment"`
	}
)

var ParticipationCollection = "participations"
var ParticipationEventView = "participations_event"

func ParticipationPermission(token *vcapool.AccessToken) (err error) {
	if !(token.Roles.Validate("employee;admin") || token.PoolRoles.Validate("network;operation;education")) {
		return vcago.NewPermissionDenied(ParticipationCollection)
	}
	return
}

func ParticipationDeletePermission(token *vcapool.AccessToken) (err error) {
	if !token.Roles.Validate("employee;admin") {
		return vcago.NewPermissionDenied(ParticipationCollection)
	}
	return
}

func ParticipationUpdatePermission(token *vcapool.AccessToken, participation *Participation) (err error) {
	if !(token.Roles.Validate("employee;admin") || token.PoolRoles.Validate("network;operation;education")) && token.ID != participation.Event.EventASPID {
		return vcago.NewPermissionDenied(ParticipationCollection)
	}
	return
}

func ParticipationPipeline() (pipe *vmdb.Pipeline) {
	pipe = vmdb.NewPipeline()
	pipe.LookupUnwind(UserCollection, "user_id", "_id", "user")
	pipe.LookupUnwind(ProfileCollection, "user_id", "user_id", "user.profile")
	pipe.LookupUnwind(EventCollection, "event_id", "_id", "event")
	pipe.LookupUnwind(CrewCollection, "crew_id", "_id", "crew")
	pipe.LookupUnwind(UserCollection, "event.event_asp_id", "_id", "event.event_asp")
	pipe.LookupUnwind(ProfileCollection, "event.event_asp_id", "user_id", "event.event_asp.profile")
	pipe.LookupUnwind(UserCollection, "event.internal_asp_id", "_id", "event.internal_asp")
	pipe.LookupUnwind(ProfileCollection, "event.internal_asp_id", "user_id", "event.internal_asp.profile")
	pipe.LookupUnwind(UserCollection, "event.creator_id", "_id", "event.creator")
	pipe.LookupUnwind(ProfileCollection, "event.creator_id", "user_id", "event.creator.profile")
	pipe.Lookup(ArtistCollection, "event.artist_ids", "_id", "event.artists")
	pipe.LookupUnwind(OrganizerCollection, "event.organizer_id", "_id", "event.organizer")
	pipe.LookupUnwind(CrewCollection, "event.crew_id", "_id", "event.crew")
	return
}

func ParticipationAspPipeline() (pipe *vmdb.Pipeline) {
	pipe = vmdb.NewPipeline()
	pipe.LookupUnwind(EventCollection, "event_id", "_id", "event")
	pipe.LookupUnwind(UserCollection, "event.event_asp_id", "_id", "event.event_asp")
	pipe.LookupUnwind(ProfileCollection, "event.event_asp_id", "user_id", "event.event_asp.profile")
	return
}

func (i *ParticipationCreate) ParticipationDatabase(token *vcapool.AccessToken) *ParticipationDatabase {
	return &ParticipationDatabase{
		ID:       uuid.NewString(),
		UserID:   token.ID,
		EventID:  i.EventID,
		Comment:  i.Comment,
		Status:   "requested",
		CrewID:   token.CrewID,
		Modified: vmod.NewModified(),
	}
}

func (i *ParticipationImport) ParticipationDatabase() *ParticipationDatabase {
	return &ParticipationDatabase{
		ID:       uuid.NewString(),
		Comment:  i.Comment,
		Status:   "requested",
		Modified: vmod.NewModified(),
	}
}
func (i *ParticipationStateRequest) IsRequested() bool {
	if i.Status == "requested" {
		return true
	}
	return false
}
func (i *ParticipationStateRequest) IsConfirmed() bool {
	if i.Status == "confirmed" {
		return true
	}
	return false
}
func (i *ParticipationStateRequest) IsWithdrawn() bool {
	if i.Status == "withdrawn" {
		return true
	}
	return false
}
func (i *ParticipationStateRequest) IsRejected() bool {
	if i.Status == "rejected" {
		return true
	}
	return false
}

func (i *ParticipationQuery) PermittedFilter(token *vcapool.AccessToken) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualStringList("_id", i.ID)
	filter.EqualStringList("event_id", i.EventID)
	filter.EqualStringList("status", i.Status)
	filter.EqualStringList("comment", i.Comment)
	filter.EqualStringList("user_id", i.UserId)
	filter.EqualStringList("crew.name", i.CrewName)
	if !token.Roles.Validate("employee;admin") {
		filter.EqualString("crew_id", token.CrewID)
	} else {
		filter.EqualStringList("crew_id", i.CrewId)
	}
	return filter.Bson()
}

func (i *ParticipationQuery) FilterUser(token *vcapool.AccessToken) *vmdb.Filter {
	filter := vmdb.NewFilter()
	filter.EqualStringList("_id", i.ID)
	filter.EqualStringList("event_id", i.EventID)
	filter.EqualStringList("status", i.Status)
	filter.EqualStringList("comment", i.Comment)
	filter.EqualStringList("crew.name", i.CrewName)
	filter.EqualStringList("crew_id", i.CrewId)
	filter.EqualString("user_id", token.ID)
	return filter
}

func (i *ParticipationQuery) FilterAspInformation(token *vcapool.AccessToken) *vmdb.Filter {
	filter := vmdb.NewFilter()
	filter.EqualStringList("event_id", i.EventID)
	if !token.Roles.Validate("employee;admin") {
		filter.EqualString("status", "confirmed")
		filter.EqualString("user_id", token.ID)
	}
	return filter
}

func (i *EventParam) FilterEvent(token *vcapool.AccessToken) *vmdb.Filter {
	filter := vmdb.NewFilter()
	filter.EqualString("event_id", i.ID)
	if !(token.Roles.Validate("employee;admin") || token.PoolRoles.Validate("network;operation;education")) {
		filter.EqualString("event.event_asp_id", token.ID)
	} else if !token.Roles.Validate("employee;admin") {
		filter.EqualString("event.crew_id", token.CrewID)
	}
	return filter
}

func (i *ParticipationDatabase) Match() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	return filter.Bson()
}

func (i *ParticipationUpdate) Match() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	return filter.Bson()
}

func (i *ParticipationStateRequest) Match() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	return filter.Bson()
}

func (i *ParticipationUpdate) PermittedFilter(token *vcapool.AccessToken) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	if token.PoolRoles.Validate("network;operation;education") {
		filter.EqualString("crew_id", token.CrewID)
	}
	return filter.Bson()
}

func (i *ParticipationParam) PermittedFilter(token *vcapool.AccessToken) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	if !(token.Roles.Validate("employee;admin") || token.PoolRoles.Validate("network;operation;education")) {
		filter.EqualString("user_id", token.ID)
	} else if !token.Roles.Validate("employee;admin") {
		filter.EqualString("crew_id", token.CrewID)
	}
	return filter.Bson()
}

func (i *ParticipationStateRequest) PermittedFilter(token *vcapool.AccessToken) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	if i.IsWithdrawn() {
		filter.EqualString("user_id", token.ID)
	} else if i.IsConfirmed() || i.IsRejected() {
		if !token.Roles.Validate("employee;admin") {
			filter.EqualString("crew_id", token.CrewID)
		}
	} else {
		filter.EqualString("_id", "not_defined")
	}
	return filter.Bson()
}
