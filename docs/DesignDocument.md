# Resume MCP Design Document

- Record work experience in DB (Sqlite locally)
- Use AI Agent to generate resume using template engine on demand
- Generate Previewable pdf
- Build on top of Go
- MCP
- Open the pdf automatically (Generated in memory and serve through url)

## Technical stacks

- https://github.com/mark3labs/mcp-go for creating mcp server
- Fiber for http endpoint
- gorm with sqlite for data storage
- go template for template language

## MCP tools:

- Create new resume with user’s name, photo and description
- Update basic info
- Add contact info with key and value
- Add work experience (start, end time, company name and job title)
- Add education experience (start, end time, and school name)
- Add feature map to the work/edu experience like

```json
{
  "gpa": 4.0,
  "features": [],
  "salary": "500HKD/month"
}
```

the tool takes three args (experience id, key, and value (any))

- Update feature map by id, key, with new value
- Delete feature map by id
- Add other experience with category of the experience and it also supports feature maps
- Get resume by name: this tool will return a structure like this

```json
{
  "resume": [
    {
      "name": "basic-info",
      "data": [
        {
          "name": "contact",
          "data": {
            "tel": "123"
          }
        }
      ]
    }
  ]
}
```

- list all saved resume and return name and id. This can be use for get resume’s detail
- delete resume by id
- generate preview pdf: this tool will return a preview link for the resume rendered by Go template engine and output a pdf. It will also open the browser to that website as well

## Rest API

The restful api will only provide a preview functionality that when mcp tool generate a preview url, it will render it as pdf file.

## Workflow

The workflow for this task is:

1.  When ai calls generate_resume function, it needs to pass a template using go template language
2.  Our server will store this info to the db and return a preview url like: https://localhost:8080/resume/preview/:sid. Then agent will return this info back to user
3.  When user goes to this url, our server will go to db to fetch the template and render the pdf on the fly
