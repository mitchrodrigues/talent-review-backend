package reviews

import (
	"fmt"
	"time"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/golly-go/plugins/gql"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/common"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/reviews/feedback"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/helpers"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/identity"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/pagination"
	"gorm.io/gorm"
)

var (
	feedbackType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Feedback",
		Fields: graphql.Fields{
			"id": {
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(Feedback).ID, nil
				},
			},
			"createdAt": {
				Type: graphql.DateTime,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(Feedback).CreatedAt, nil
				},
			},
			"collectionEndAt": {
				Type: graphql.DateTime,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(Feedback).CollectionEndAt, nil
				},
			},
			"email": {
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(Feedback).Email, nil
				},
			},
			"submittedAt": {
				Type: graphql.DateTime,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if t := p.Source.(Feedback).SubmittedAt; t.IsZero() {
						return nil, nil
					}

					return p.Source.(Feedback).SubmittedAt, nil
				},
			},
			"employee": {
				Type: employees.EmployeeSimplified,
				Resolve: gql.NewHandler(gql.Options{
					Public: true,
					Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
						employeeID := params.Source.(Feedback).EmployeeID

						return golly.LoadData(
							ctx.Context,
							fmt.Sprintf("employee:%s", employeeID),
							func(gctx golly.Context) (employees.Employee, error) {
								return employees.FindEmployeeByID_Unsafe(ctx.Context, employeeID)
							},
						)
					},
				}),
			},
		},
	})

	query = graphql.Fields{
		//********** Feedback ***************//
		"feedbacks": {
			Name: "feedbacks",
			Args: graphql.FieldConfigArgument{
				"pagination": &graphql.ArgumentConfig{
					Type: pagination.PaginationInputType,
				},
			},
			Type: pagination.PaginationType[Feedback](feedbackType),
			Resolve: gql.NewHandler(gql.Options{
				Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
					ident := identity.FromContext(ctx.Context)

					employeeIDs, err := employees.FindEmployeeIDsByManagerUserID(ctx.Context, ident.UID)
					if err != nil {
						return nil, err
					}

					return pagination.
						NewCursorPaginationFromArgs(
							params.Args,
							[]Feedback{},
							common.OrganizationIDScopeForContext(ctx.Context),
							func(db *gorm.DB) *gorm.DB {
								return db.Where("employee_id IN ?", employeeIDs)
							},
						).
						Paginate(ctx.Context)
				},
			}),
		},
		"feedbackForCode": {
			Name: "feedbackForCode",
			Args: graphql.FieldConfigArgument{
				"code": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Type: feedbackType,
			Resolve: gql.NewHandler(gql.Options{
				Public: true,
				Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
					feedback, err := FindFeedbackForCode(ctx.Context, params.Args["code"].(string))
					if feedback.ID == uuid.Nil {
						return nil, err
					}

					return feedback, nil
				},
			}),
		},
	}

	// type CreateBulkFeedbackInput struct {
	// 	EmployeeIDs      []uuid.UUID
	// 	AdditionalEmails []string
	// 	IncludeTeam      bool
	// 	CollectionEndAt  time.Time
	// }
	createTeamInputType = graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "CreateFeedbacksInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"employeeIDs":      {Type: graphql.NewNonNull(graphql.NewList(graphql.String))},
			"additionalEmails": {Type: graphql.NewList(graphql.String)},
			"includeTeam":      {Type: graphql.Boolean},
			"collectionEndAt":  {Type: graphql.NewNonNull(graphql.DateTime)},
		},
	})

	udpateFeedbackDetailsType = graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "UpdateFeedbackDetailsInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"strengths":     {Type: graphql.String},
			"opportunities": {Type: graphql.String},
			"additional":    {Type: graphql.String},
			"rating":        {Type: graphql.Int},
			"enoughData":    {Type: graphql.Boolean},
		},
	})

	mutations = graphql.Fields{
		//********** Feedback ***************//
		"submitFeedback": {
			Name: "submitFeedback",
			Type: feedbackType,
			Args: graphql.FieldConfigArgument{
				"code": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"id":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: gql.NewHandler(gql.Options{
				Public: true,
				Handler: func(wctx golly.WebContext, params gql.Params) (interface{}, error) {
					id, err := helpers.ExtractAndParseUUID(params.Args, "id")
					if err != nil {
						return nil, err
					}

					fb, err := FindFeedbackByIDAndCode(wctx.Context, id, params.Args["code"].(string))
					if err != nil {
						return nil, err
					}

					_, gctx := identity.SetOrganizationID(wctx.Context, fb.OrganizationID)

					err = eventsource.Call(gctx, &fb.Aggregate, feedback.Submit{}, params.Metadata())
					return fb, err
				},
			}),
		},

		"updateFeedbackDetails": {
			Name: "updateFeedbackDetails",
			Type: feedbackType,
			Args: graphql.FieldConfigArgument{
				"code":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(udpateFeedbackDetailsType)},
			},
			Resolve: gql.NewHandler(gql.Options{
				Public: true,
				Handler: func(wctx golly.WebContext, params gql.Params) (interface{}, error) {
					id, err := helpers.ExtractAndParseUUID(params.Args, "id")
					if err != nil {
						return nil, err
					}

					fb, err := FindFeedbackByIDAndCode(wctx.Context, id, params.Args["code"].(string))
					if err != nil {
						return nil, err
					}

					var strength string = ""
					if val, err := helpers.ExtractArg[string](params.Input, "strengths"); err == nil {
						strength = val
					} else {
						wctx.Logger().Warnf("Strengths: %#v", err)
					}

					var opportunities string = ""
					if val, err := helpers.ExtractArg[string](params.Input, "opportunities"); err == nil {
						opportunities = val
					}

					var additional string = ""
					if val, err := helpers.ExtractArg[string](params.Input, "additional"); err == nil {
						additional = val
					}

					var rating int = 0
					if val, err := helpers.ExtractArg[int64](params.Input, "rating"); err == nil {
						rating = int(val)
					}

					var enoughData *bool = nil
					if val, err := helpers.ExtractArg[bool](params.Input, "enoughData"); err == nil {
						enoughData = &val
					}

					_, gctx := identity.SetOrganizationID(wctx.Context, fb.OrganizationID)

					err = eventsource.Call(gctx, &fb.Aggregate, feedback.CreateOrUpdateDetails{
						Strength:      strength,
						Opportunities: opportunities,
						Additional:    additional,
						Rating:        rating,
						EnoughData:    enoughData,
					}, params.Metadata())

					return fb, err
				},
			}),
		},

		"createFeedbacks": {
			Name: "createFeedbacks",
			Type: graphql.NewList(feedbackType),
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(createTeamInputType)},
			},

			Resolve: gql.NewHandler(gql.Options{
				Handler: func(ctx golly.WebContext, params gql.Params) (interface{}, error) {
					employeeIDs := golly.Map(params.Input["employeeIDs"].([]interface{}), func(id interface{}) uuid.UUID {
						return uuid.MustParse(id.(string))
					})

					includeTeam, _ := helpers.ExtractArg[bool](params.Input, "includeTeam")
					additionalEmails, _ := helpers.ExtractArg[[]interface{}](params.Input, "additionalEmails")

					return CreateBulkFeedback(ctx.Context, CreateBulkFeedbackInput{
						EmployeeIDs:     employeeIDs,
						IncludeTeam:     includeTeam,
						CollectionEndAt: params.Input["collectionEndAt"].(time.Time),
						AdditionalEmails: golly.Map(additionalEmails, func(i interface{}) string {
							return i.(string)
						}),
					}, params.Metadata())

				},
			}),
		},
	}
)

func InitGraphQL() {
	gql.RegisterQuery(query)
	gql.RegisterMutation(mutations)
}
